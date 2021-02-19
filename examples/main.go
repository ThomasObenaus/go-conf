package main

import (
	"fmt"
	"reflect"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/ThomasObenaus/go-conf/interfaces"

	//"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
)

// TODO: Fail in case there are duplicate settings (names) configured
// TODO: Check if pointer fields are supported
// TODO: Add support for shorthand flags
// TODO: Think about required vs. optional (explicit vs implicit)

type Cfg struct {
	DryRun        bool           // this should be ignored since its not annotated, but it can be still read using on the usual way
	Name          string         `cfg:"{'name':'name','desc':'the name of the config','short':'n'}"`
	Prio          int            `cfg:"{'name':'prio','desc':'the prio'}"`
	Immutable     bool           `cfg:"{'name':'immutable','desc':'can be modified or not','default':false}"`
	NumericLevels []int          `cfg:"{'name':'numeric-levels','desc':'allowed levels','default':[1,2]}"`
	Levels        []string       `cfg:"{'name':'levels','desc':'allowed levels','default':['a','b']}"`
	ConfigStore   configStore    `cfg:"{'name':'config-store','desc':'the config store'}"`
	TargetSecrets []targetSecret `cfg:"{'name':'target-secrets','desc':'list of target secrets','default':[{'name':'1mysecret1','key':'sdlfks','count':231},{'name':'mysecret2','key':'sdlfks','count':231}]}"`
}

type configStore struct {
	FilePath     string        `cfg:"{'name':'file-path','desc':'the path','default':'configs/'}"`
	TargetSecret targetSecret  `cfg:"{'name':'target-secret','desc':'the secret'}"`
	LogLevel     zerolog.Level `cfg:"{'name':'log-level','default':'info','mapfun':'strToLogLevel'}"`
}

type targetSecret struct {
	Name  string `cfg:"{'name':'name','desc':'the name of the secret'}"`
	Key   string `cfg:"{'name':'key','desc':'the key'}"`
	Count int    `cfg:"{'name':'count','desc':'the count','default':0}"`
}

func main() {

	args := []string{
		"--dry-run",
		//"--name=hello",
		"-n=hello",
		"--prio=23",
		"--immutable=true",
		"--numeric-levels=1,2,3",
		"--config-store.file-path=/devops",
		"--config-store.target-secret.key=#lsdpo93",
		"--config-store.target-secret.name=mysecret",
		"--config-store.target-secret.count=2323",
		"--config-store.log-level=fatal",
		"--target-secrets=[{'name':'mysecret1','key':'sdlfks','count':231},{'name':'mysecret2','key':'sdlfks','count':231}]",
	}

	parsedConfig, err := New(args, "ABCDE")
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Println("SUCCESS")
	fmt.Printf("%v", parsedConfig)
	//spew.Dump(parsedConfig)
}

var dryRun = config.NewEntry("dry-run", "If true, then sokar won't execute the planned scaling action. Only scaling\n"+
	"actions triggered via ScaleBy end-point will be executed.", config.Default(false))
var configEntries = []config.Entry{
	dryRun,
}

func New(args []string, serviceAbbreviation string) (Cfg, error) {
	cfg := Cfg{}

	provider, err := config.NewConfigProvider(
		&cfg,
		serviceAbbreviation,
		serviceAbbreviation,
		config.CustomConfigEntries(configEntries),
		config.Logger(interfaces.DebugLogger),
	)
	if err != nil {
		return Cfg{}, err
	}

	provider.RegisterMappingFunc("strToLogLevel", func(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
		asStr, ok := rawUntypedValue.(string)
		if !ok {
			return nil, fmt.Errorf("Expected a string. Type '%T' is not supported", rawUntypedValue)
		}
		return zerolog.ParseLevel(asStr)
	})

	err = provider.ReadConfig(args)
	if err != nil {

		fmt.Print(provider.Usage())

		return Cfg{}, err
	}

	if err := cfg.overWriteCfgValues(provider); err != nil {
		return Cfg{}, err
	}

	return cfg, nil
}

func (cfg *Cfg) overWriteCfgValues(provider interfaces.Provider) error {
	cfg.DryRun = provider.GetBool(dryRun.Name())
	cfg.Name = "Thomas (OVERWRITTEN)"
	return nil
}
