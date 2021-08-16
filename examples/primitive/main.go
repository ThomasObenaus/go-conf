package main

import (
	"fmt"
	"os"
	"time"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

// PrimitiveTypes Define the config struct and annotate it with the cfg tag.
type PrimitiveTypes struct {
	Field1 string        `cfg:"{'name':'field-1'}"`
	Field2 int           `cfg:"{'name':'field-2'}"`
	Field3 float64       `cfg:"{'name':'field-3'}"`
	Field4 bool          `cfg:"{'name':'field-4'}"`
	Field5 time.Duration `cfg:"{'name':'field-5'}"`
}

func main() {

	args := []string{
		"--field-1=one",
		"--field-2=1234",
		"--field-3=12.34",
		"--field-4=true",
		"--field-5=1234m",
	}

	// 1. Create an instance of the config struct that should be filled
	cfg := PrimitiveTypes{}

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

	// 3. Start reading and fill the config parameters
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
