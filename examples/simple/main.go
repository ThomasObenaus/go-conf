package main

import (
	"fmt"
	"os"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

// MyFontConfig an example configuration.
//
// Just a struct with the desired config parameters has to be defined.
// Each field has to be annotated with the cfg struct tag
//
//	`cfg:{'name':'<name of the parameter>','desc':'<description>'}`
//
// The tag has to be specified as json structure using single quotes.
// Mandatory fields are 'name' and 'desc'.
type MyFontConfig struct {
	Color string `cfg:"{'name':'color','desc':'The value of the color as hexadecimal RGB string.','default':'#FFFFFF'}"`
	Name  string `cfg:"{'name':'name','desc':'Name of the font to be used.'}"`
	Size  int    `cfg:"{'name':'size','desc':'Size of the font.','short':'s'}"`
}

func main() {

	args := []string{
		// color not set  --> default value will be used "--color=#ff00ff",
		"--name=Arial",
		"-s=12", // use -s (short hand version) instead of --size
	}

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
		os.Exit(1)
	}

	fmt.Println("##### Successfully read the config")
	fmt.Println()
	spew.Dump(cfg)
}
