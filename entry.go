package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Entry is one item to define a configuration
type Entry struct {
	name         string
	usage        string
	defaultValue interface{}
	desiredType  reflect.Type

	bindFlag      bool
	flagShortName string

	bindEnv bool
}

// EntryOption represents an option for the Entry
type EntryOption func(e *Entry)

// Default specifies a default value
func Default(value interface{}) EntryOption {
	return func(e *Entry) {
		e.defaultValue = value
	}
}

// ShortName specifies the shorthand (one-letter) flag name
func ShortName(fShort string) EntryOption {
	return func(e *Entry) {
		e.flagShortName = fShort
	}
}

// Bind enables/ disables binding of flag and env var
func Bind(flag, env bool) EntryOption {
	return func(e *Entry) {
		e.bindFlag = flag
		e.bindEnv = env
	}
}

// DesiredType sets the desired type of this entry
func DesiredType(t reflect.Type) EntryOption {
	return func(e *Entry) {
		e.desiredType = t
	}
}

// NewEntry creates a new Entry that is available as flag, config file entry and environment variable
func NewEntry(name, usage string, options ...EntryOption) Entry {
	entry := Entry{
		name:          name,
		usage:         usage,
		flagShortName: "",
		defaultValue:  nil,
		bindFlag:      true,
		bindEnv:       true,
		desiredType:   nil,
	}

	// apply the options
	for _, opt := range options {
		opt(&entry)
	}

	// try the best to deduce the desired type
	if entry.desiredType == nil {
		if entry.defaultValue != nil {
			entry.desiredType = reflect.TypeOf(entry.defaultValue)
		} else {
			entry.desiredType = reflect.TypeOf("")
		}
	}

	return entry
}

func (e Entry) String() string {
	return fmt.Sprintf("--%s (-%s) [default:%v (%T)]\t- %s", e.name, e.flagShortName, e.defaultValue, e.defaultValue, e.usage)
}

// Name provides the specified name for this entry
func (e Entry) Name() string {
	return e.name
}

// IsRequired returns true in case no default value is given
func (e Entry) IsRequired() bool {
	return e.defaultValue == nil
}

func checkViper(vp *viper.Viper, entry Entry) error {
	if vp == nil {
		return fmt.Errorf("Viper is nil")
	}

	if len(entry.name) == 0 {
		return fmt.Errorf("Name is missing")
	}

	return nil
}

func registerFlag(flagSet *pflag.FlagSet, entry Entry) error {
	if !entry.bindFlag {
		return nil
	}
	if flagSet == nil {
		return fmt.Errorf("FlagSet is nil")
	}
	if len(entry.name) == 0 {
		return fmt.Errorf("Name is missing")
	}

	valueDesiredType := entry.defaultValue

	if valueDesiredType == nil {
		valueType := reflect.New(entry.desiredType)
		if valueType.Kind() != reflect.Ptr {
			return fmt.Errorf("Failed deducing desired type for entry %v", entry)
		}
		valueDesiredType = valueType.Elem().Interface()
	}

	// TODO: Regard default value and set it when registering the flag

	switch castedDefaultValue := valueDesiredType.(type) {
	case string:
		flagSet.StringP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case uint:
		flagSet.UintP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case int:
		flagSet.IntP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case bool:
		flagSet.BoolP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case time.Duration:
		flagSet.DurationP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case float32:
		flagSet.Float32P(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case float64:
		flagSet.Float64P(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []bool:
		flagSet.BoolSliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []string:
		flagSet.StringSliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []time.Duration:
		flagSet.DurationSliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []int:
		flagSet.IntSliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []int32:
		flagSet.Int32SliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []int64:
		flagSet.Int64SliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []uint:
		flagSet.UintSliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []float64:
		flagSet.Float64SliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	case []float32:
		flagSet.Float32SliceP(entry.name, entry.flagShortName, castedDefaultValue, entry.usage)
	default:
		typeOfDefaultValue := reflect.TypeOf(castedDefaultValue)
		if typeOfDefaultValue.Kind() != reflect.Slice {
			return fmt.Errorf("Type %s is not yet supported", typeOfDefaultValue)
		}

		// this part supports slices of custom structs and registers the according flag for it
		s := sliceOfMapStringToInterfaceFlag{}
		flagSet.VarP(&s, entry.name, entry.flagShortName, entry.usage)
		return nil
	}

	return nil
}

func setDefault(vp *viper.Viper, entry Entry) error {
	if err := checkViper(vp, entry); err != nil {
		return err
	}

	if entry.defaultValue != nil {
		vp.SetDefault(entry.name, entry.defaultValue)
	}

	return nil
}

func registerEnv(vp *viper.Viper, envPrefix string, entry Entry) error {
	if !entry.bindEnv {
		return nil
	}
	if err := checkViper(vp, entry); err != nil {
		return err
	}

	if len(envPrefix) > 0 {
		vp.SetEnvPrefix(envPrefix)
	}
	return vp.BindEnv(entry.name)
}

// sliceOfMapStringToInterfaceFlag is a struct that can be used to represent a flag of type
//	[]map[string]interface{}
// That is a slice of arbitrary structs.
type sliceOfMapStringToInterfaceFlag struct {
	// used to be returned in the String method. Its better to return the value as
	// json string since this can be parsed easier if needed.
	jsonSting string
}

func (l *sliceOfMapStringToInterfaceFlag) String() string {
	return l.jsonSting
}

func (l *sliceOfMapStringToInterfaceFlag) Type() string {
	return "[]map[string]interface{}"
}

func (l *sliceOfMapStringToInterfaceFlag) Set(in string) error {
	in = strings.ReplaceAll(in, "'", "\"")
	l.jsonSting = in
	return nil
}
