package main

import (
	"fmt"
	"os"

	config "github.com/ThomasObenaus/go-conf"
)

type MyFontConfig struct {
	Color string `cfg:"{'name':'color','desc':'A required parameter (no default is defined)'}"`
	Size  int    `cfg:"{'name':'size','desc':'An optional parameter (since a default value is defined)','default':12}"`
}

func main() {

	args := []string{}

	// 1. Create an instance of the config struct that should be filled
	cfg := MyFontConfig{}

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

		// It is actually expected to fail here (since the required parameter is not provided via args)
		// Hence we return with 0
		os.Exit(0)
	}

	// It is not expected to come here. Hence we return with != 0
	os.Exit(1)
}
