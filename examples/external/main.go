package main

import (
	"fmt"
	"os"
	"time"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/ThomasObenaus/go-conf/interfaces"
	"github.com/davecgh/go-spew/spew"
)

type PrimitiveTypes struct {
	Field1 time.Time `cfg:"{'name':'field-1'}"`
}

func main() {

	args := []string{
		"--field-1=22", //2021-02-22 12:34:00 +0000 UTC",
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
		config.Logger(interfaces.DebugLogger),
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
