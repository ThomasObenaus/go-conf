package main

import (
	"fmt"
	"os"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

// MyConfig Define the config struct and annotate it with the cfg tag.
type MyConfig struct {
	// --default-only not set in file, environment and cli --> will use the default value
	DefaultOnly string `cfg:"{'name':'default-only','default':'Default'}"`
	// --from-file only set in in file, not in environment and cli --> will use the value from file
	FromFile string `cfg:"{'name':'from-file','default':'Default'}"`
	// --from-env set in file and environment, but not in cli --> will use the env value
	FromEnv string `cfg:"{'name':'from-env','default':'Default'}"`
	// --from-cli set in file, environment and cli --> will use the cli value
	FromCLI string `cfg:"{'name':'from-cli','default':'Default'}"`
}

func main() {

	os.Setenv("MY_APP_FROM_ENV", "env")
	args := []string{
		"--config-file-name=examples/multisource/config.yaml",
		"--from-cli=cli",
	}

	// 1. Create an instance of the config struct that should be filled
	cfg := MyConfig{}

	// 2. Create an instance of the config provider which is responsible to read the config
	// from defaults, environment variables, config file or command line
	prefixForEnvironmentVariables := "MY_APP"
	nameOfTheConfig := "MY_APP"
	provider, err := config.NewConfigProvider(
		&cfg,
		nameOfTheConfig,
		prefixForEnvironmentVariables,
		config.CfgFile("config-file-name", "f"), // Overwrites the default parameter (--config-file) for the config file
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
