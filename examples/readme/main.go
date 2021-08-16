package main

import (
	"fmt"

	config "github.com/ThomasObenaus/go-conf"
)

// MyFontConfig Define the config struct and annotate it with the cfg tag.
type MyFontConfig struct {
	Color string `cfg:"{'name':'color','desc':'The value of the color as hexadecimal RGB string.','default':'#FFFFFF'}"`
	Name  string `cfg:"{'name':'name','desc':'Name of the font to be used.'}"`
	Size  int    `cfg:"{'name':'size','desc':'Size of the font.','short':'s'}"`
}

func main() {

	// Some command line arguments
	args := []string{
		// color not set  --> default value will be used "--color=#ff00ff",
		"--name=Arial",
		"-s=12", // use -s (short hand version) instead of --size
	}

	// 1. Create an instance of the config struct that should be filled
	cfg := MyFontConfig{}

	// 2. Create an instance of the config provider
	provider, err := config.NewConfigProvider(&cfg, "MY_APP", "MY_APP")
	if err != nil {
		panic(err)
	}

	// 3. Read the config and populate the struct
	if err := provider.ReadConfig(args); err != nil {
		panic(err)
	}

	// 4. Thats it! Now the config can be used.
	fmt.Printf("FontConfig: color=%s, name=%s, size=%d\n", cfg.Color, cfg.Name, cfg.Size)
}
