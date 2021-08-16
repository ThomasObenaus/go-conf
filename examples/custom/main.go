package main

import (
	"fmt"
	"os"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

// MyConfig Define the config struct and annotate it with the cfg tag.
type MyConfig struct {
	Field1 string `cfg:"{'name':'field-1'}"`
	Field2 bool   // not annotated, since it is filled manually
	Field3 int    // not annotated, since it is filled manually
}

func main() {

	args := []string{
		"--field-1=one",
		"--field-2=true",
		"--field-3=5678",
	}

	// 1. Create custom config entries
	customEntries := []config.Entry{
		config.NewEntry("field-2", "A bool flag", config.Default(false)),
		config.NewEntry("field-3", "A integer flag", config.Default(1234)),
	}

	// 2. Create an instance of the config struct that should be filled
	cfg := MyConfig{}

	// 3. Create an instance of the config provider which is responsible to read the config
	// from defaults, environment variables, config file or command line
	prefixForEnvironmentVariables := "MY_APP"
	nameOfTheConfig := "MY_APP"
	provider, err := config.NewConfigProvider(
		&cfg,
		nameOfTheConfig,
		prefixForEnvironmentVariables,
		config.CustomConfigEntries(customEntries), // 4. register the custom config entries
	)
	if err != nil {
		panic(err)
	}

	// 5. Start reading and fill the config parameters
	err = provider.ReadConfig(args)
	if err != nil {
		fmt.Println("##### Failed reading the config")
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Println("Usage:")
		fmt.Print(provider.Usage())
		os.Exit(1)
	}

	// 6. Manually fill the struct reading the custom config entries
	cfg.Field2 = provider.GetBool(customEntries[0].Name())
	cfg.Field3 = provider.GetInt(customEntries[1].Name())

	fmt.Println("##### Successfully read the config")
	fmt.Println()
	spew.Dump(cfg)
}
