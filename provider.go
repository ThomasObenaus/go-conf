package config

import (
	"fmt"
	"strings"

	"github.com/ThomasObenaus/go-conf/interfaces"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// providerImpl is a structure containing the parsed configuration
type providerImpl struct {
	// config entries are all definitions of config entries that should be regarded
	configEntries []Entry

	// this special entry is used to specify the name + location of the config file
	configFileEntry Entry

	configName string

	// the environment prefix (will be added to all env vars) <envPrefix>_<name of config entry>
	// e.g. assuming the envPrefix is "myApp" and the name of the config entry is "my-entry"
	// then the env var is MYAPP_MY_ENTRY
	envPrefix string

	// instance of pflag, needed to parse command line parameters
	pFlagSet *pflag.FlagSet

	// instance of viper, needed to parse env vars and to read from cfg-file
	*viper.Viper

	// configTarget is the (user) configuration object the configuration should be applied to.
	// It can be nil, hence the code has to cover this case.
	configTarget interface{}

	logger interfaces.LoggerFunc

	mappingFuncRegistry map[string]interfaces.MappingFunc
}

// ProviderOption represents an option for the Provider
type ProviderOption func(p *providerImpl)

// CfgFile specifies a default value
func CfgFile(parameterName, shortParameterName string) ProviderOption {
	return func(p *providerImpl) {
		p.configFileEntry = NewEntry(parameterName, "Specifies the full path and name of the configuration file", ShortName(shortParameterName))
	}
}

// CustomConfigEntries allows to add config entries that are created manually via NewEntry(..)
func CustomConfigEntries(customConfigEntries []Entry) ProviderOption {
	return func(p *providerImpl) {
		if p.configEntries == nil {
			p.configEntries = make([]Entry, 0)
		}
		p.configEntries = append(p.configEntries, customConfigEntries...)
	}
}

// Logger exchanges the logger function. This provides the possibility to integrate your own logger.
// Per default the NoLogging function is used (disables logging completely).
// Other predefined logger functions (based on fmt.Printf) are DebugLogger, InfoLogger, WarnLogger and ErrorLogger.
func Logger(logger interfaces.LoggerFunc) ProviderOption {
	return func(p *providerImpl) {
		p.logger = logger
	}
}

// NewProvider creates a new config provider that is able to parse the command line, env vars and config file based
// on the given entries.
//
// DEPRECATED: Please use NewConfigProvider instead.
func NewProvider(configEntries []Entry, configName, envPrefix string, options ...ProviderOption) (interfaces.Provider, error) {
	opt := CustomConfigEntries(configEntries)
	options = append(options, opt)
	provider, err := NewConfigProvider(nil, configName, envPrefix, options...)
	if err != nil {
		return nil, err
	}

	provider.Log(interfaces.LogLevel_Warn, "You are using the old, deprecated config interface 'NewProvider' please use 'NewConfigProvider' instead.")
	return provider, nil
}

// NewConfigProvider creates a new config provider that is able to parse the command line, env vars and config file based
// on the given entries. This config provider automatically generates the needed config entries and fills the given config target
// based on the annotations on this struct.
// In case custom config entries should be used beside the annotations on the struct one can define them via
//	CustomConfigEntries(customEntries)`
// e.g.
//
//	customEntries:=[]Entry{
//	// fill entries here
//	}
//	provider,err := NewConfigProvider(&myConfig,"my-config","MY_APP",CustomConfigEntries(customEntries))
func NewConfigProvider(target interface{}, configName, envPrefix string, options ...ProviderOption) (interfaces.Provider, error) {
	defaultConfigFileEntry := NewEntry("config-file", "Specifies the full path and name of the configuration file", Bind(true, true))
	provider := &providerImpl{
		configEntries:       make([]Entry, 0),
		configName:          configName,
		envPrefix:           envPrefix,
		pFlagSet:            pflag.NewFlagSet(configName, pflag.ContinueOnError),
		Viper:               viper.New(),
		configFileEntry:     defaultConfigFileEntry,
		configTarget:        target,
		logger:              interfaces.NoLogging,
		mappingFuncRegistry: make(map[string]interfaces.MappingFunc),
	}

	// apply the options
	for _, opt := range options {
		opt(provider)
	}

	if provider.logger == nil {
		return nil, fmt.Errorf("The Logger set via config.Logger must not be nil")
	}

	// Enable casting to type based on given default values
	// this ensures that viper.Get() returns the casted instance instead of the plain value.
	// That helps for example when a configuration is of type time.Duration.
	// Usually viper.Get() would return a string but now it returns a time.Duration
	provider.Viper.SetTypeByDefaultValue(true)

	// For backwards compatibility we also allow to provide no target (this will be the case if the NewProvider function is used)
	if provider.configTarget != nil {
		configEntries, err := CreateEntriesFromStruct(provider.configTarget, provider.Log)
		if err != nil {
			return nil, errors.Wrapf(err, "Extracting configuration annotations")
		}
		provider.configEntries = append(provider.configEntries, configEntries...)
	} else {
		provider.logger(interfaces.LogLevel_Info, "No target given. Hence the config is not automatically processed and applied.")
	}

	return provider, nil
}

func (p *providerImpl) Log(lvl interfaces.LogLevel, formatString string, a ...interface{}) {
	p.logger(lvl, formatString, a...)
}

func (p *providerImpl) String() string {
	return fmt.Sprintf("%s: %v", p.configName, p.AllSettings())
}

func (p *providerImpl) RegisterMappingFunc(name string, mFunc interfaces.MappingFunc) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return fmt.Errorf("Name for the mapping function must be provided")
	}

	if mFunc == nil {
		return fmt.Errorf("Mapping function must not be nil")
	}

	if old, ok := p.mappingFuncRegistry[name]; ok {
		p.logger(interfaces.LogLevel_Warn, "Overwriting mapping function '%v' with '%v' because both were registered with the same name '%s'", old, mFunc, name)
	}

	p.mappingFuncRegistry[name] = mFunc

	return nil
}

func (p *providerImpl) Usage() string {
	entriesAsString := make([]string, 0)

	for _, entry := range p.configEntries {
		entryDefinition := entryDefinitionAsString(entry)
		entriesAsString = append(entriesAsString, entryDefinition)

		// default
		defaultStr := "n/a"
		if entry.defaultValue != nil {
			defaultStr = fmt.Sprintf("%v (type=%T)", entry.defaultValue, entry.defaultValue)
		}
		entriesAsString = append(entriesAsString, fmt.Sprintf("\tdefault: %s", defaultStr))

		// usage
		usageStr := wrapText(fmt.Sprintf("desc: %s", entry.usage), 140, "\n\t")
		entriesAsString = append(entriesAsString, fmt.Sprintf("\t%s", usageStr))
		entriesAsString = append(entriesAsString, "")
	}
	return strings.Join(entriesAsString, "\n")
}

func entryDefinitionAsString(entry Entry) string {
	reqStr := ""
	if entry.IsRequired() {
		reqStr = " [required]"
	}
	return fmt.Sprintf("--%s (-%s)%s", entry.Name(), entry.flagShortName, reqStr)
}

func wrapText(text string, afterNChars int, wrapChar string) string {

	if len(text) <= afterNChars {
		return text
	}

	parts := []string{}
	for len(text) > afterNChars {
		text = strings.TrimSpace(text)
		parts = append(parts, text[:(afterNChars)])
		text = text[afterNChars:]
		text = strings.TrimSpace(text)
	}

	if len(text) > 0 {
		parts = append(parts, text)
	}

	return strings.Join(parts, wrapChar)
}
