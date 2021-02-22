# go-conf - Configuration with less code

[![Go Reference](https://pkg.go.dev/badge/github.com/ThomasObenaus/go-conf.svg)](https://pkg.go.dev/github.com/ThomasObenaus/go-conf) ![build](https://github.com/ThomasObenaus/go-conf/workflows/build/badge.svg?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/ThomasObenaus/go-conf)](https://goreportcard.com/report/github.com/ThomasObenaus/go-conf)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=alert_status)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=coverage)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=code_smells)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=security_rating)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ThomasObenaus_go-conf&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=ThomasObenaus_go-conf) ![CodeQL](https://github.com/ThomasObenaus/go-conf/workflows/CodeQL/badge.svg)

## Installation

```bash
go get https://github.com/ThomasObenaus/go-conf
```

## What is go-conf?

go-conf is a **solution for handling configurations** in golang applications.

go-conf **supports reading configuration parameters from multiple sources**.
The order they are applied is:

1. Default values are overwritten by
2. Parameters defined in the config-file, which are overwritten by
3. Environment variables, which are overwritten by
4. Command-Line parameters

The aim is to **write as less code as possible**:

- No need to write code to integrate multiple libraries that support reading a configuration from file/ commandline or the environment.
- No need to code to take the values from that library to fill it into the config struct you want to use in your app anyway.

Instead one just has to define the config structure and annotates it with struct tags.

```go
package main

import (
    "fmt"

    config "github.com/ThomasObenaus/go-conf"
)

// Define the config struct and annotate it with the cfg tag.
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

```

### Examples

- [examples/simple](simple): Showcases the basic features of go-conf.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FThomasObenaus%2Fgo-conf.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FThomasObenaus%2Fgo-conf?ref=badge_large)
