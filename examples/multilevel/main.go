package main

import (
	"fmt"
	"os"

	config "github.com/ThomasObenaus/go-conf"
	"github.com/davecgh/go-spew/spew"
)

// ThemeConfig configuration for a theme.
// The different levels are separated by a dot from each another when be set via command line.
// Hence
//	type ThemeConfig struct {
//		Header FormattedTextBox `cfg:"{'name':'header','desc':'The heading'}"`
//	}
//
//	type FormattedTextBox struct {
//		Value string `cfg:"{'name':'value','desc':'The content of the text box','default':''}"`
//	}
// can be set by
//	--header.value="hi there"
type ThemeConfig struct {
	Header FormattedTextBox `cfg:"{'name':'header','desc':'The heading'}"`
	Footer FormattedTextBox `cfg:"{'name':'footer','desc':'The footer'}"`
}

// FormattedTextBox used to configure the look of a text box
type FormattedTextBox struct {
	Font   Font   `cfg:"{'name':'font','desc':'Definition of the text box font.'}"`
	Border Border `cfg:"{'name':'border','desc':'Definition of the text box border.'}"`
	Value  string `cfg:"{'name':'value','desc':'The content of the text box','default':''}"`
}

// Font used to configure the look of a font
type Font struct {
	Color string `cfg:"{'name':'color','desc':'The value of the color as hexadecimal RGB string.','default':'#FFFFFF'}"`
	Name  string `cfg:"{'name':'name','desc':'Name of the font to be used.','default':'arial'}"`
	Size  int    `cfg:"{'name':'size','desc':'Size of the font.','default':12}"`
}

// Border used to configure the look of a border
type Border struct {
	Color string `cfg:"{'name':'color','desc':'The value of the color as hexadecimal RGB string.','default':'#FFFFFF'}"`
	Width int    `cfg:"{'name':'width','desc':'Width of the border.','default':1}"`
}

func main() {
	args := []string{
		"--header.font.color=#AAAAAA",
		"--header.border.width=0",
		"--header.value=This is the header",
		"--footer.font.size=10",
		"--footer.value=This is the footer",
	}

	// 1. Create an instance of the config struct that should be filled
	cfg := ThemeConfig{}

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
