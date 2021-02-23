package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

type ExternalTypes struct {
	Field1 time.Time `cfg:"{'name':'field-1','mapfun':'strToTime'}"`
}

func main() {

	args := []string{
		"--field-1=2021-02-22 12:34:00 +0000 UTC",
	}

	// 1. Create an instance of the config struct that should be filled
	cfg := ExternalTypes{}

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

	// 3. Register the mapping function (string->time.Time)
	if err := provider.RegisterMappingFunc("strToTime", strToTime); err != nil {
		panic(err)
	}

	// 4. Start reading and fill the config parameters
	err = provider.ReadConfig(args)
	if err != nil {
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

// strToTime maps a string to a time.Time
func strToTime(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
	asStr, ok := rawUntypedValue.(string)
	if !ok {
		return nil, fmt.Errorf("Expected a string. Type '%T' is not supported", rawUntypedValue)
	}
	t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", asStr)
	if err != nil {
		return nil, fmt.Errorf("Parse %s to time failed: %s", asStr, err.Error())
	}
	return t, nil
}
