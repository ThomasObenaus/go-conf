package main

import (
	"fmt"
	"os"
	"reflect"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

// LogLevel defines log levels.
type LogLevel int8

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type LogConfig struct {
	Level LogLevel `cfg:"{'name':'level','desc':'Defines the loglevel (debug|info|warn|err).','default':'info','mapfun':'strToLogLevel'}"`
}

func main() {

	args := []string{
		"--level=err",
	}

	// 1. Create an instance of the config struct that should be filled
	cfg := LogConfig{}

	// 2. Create an instance of the config provider which is responsible to read the config
	// from defaults, environment variables, config file or command line
	prefixForEnvironmentVariables := "MY_APP"
	nameOfTheConfig := "MY_APP"
	provider, err := config.NewConfigProvider(
		&cfg,
		nameOfTheConfig,
		prefixForEnvironmentVariables,
	)
	if err != nil {
		panic(err)
	}

	// 3. Register the mapping function
	if err := provider.RegisterMappingFunc("strToLogLevel", strToLogLevel); err != nil {
		panic(err)
	}

	// 4. Start reading and fill the config parameters
	if err := provider.ReadConfig(args); err != nil {
		fmt.Println("##### Failed reading the config")
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Println("Usage:")
		fmt.Print(provider.Usage())
		os.Exit(1)
	}

	fmt.Println("##### Successfully read the config")
	fmt.Println()
	spew.Dump(cfg)
}

// strToLogLevel maps a given string into a LogLevel (uint8)
func strToLogLevel(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
	// ensure that the expectation of the input type is satisfied
	asStr, ok := rawUntypedValue.(string)
	if !ok {
		return nil, fmt.Errorf("Expected a string. Type '%T' is not supported", rawUntypedValue)
	}

	// return the target type
	switch asStr {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn":
		return WarnLevel, nil
	case "err":
		return ErrorLevel, nil
	}
	return nil, fmt.Errorf("%s is unknown", asStr)
}
