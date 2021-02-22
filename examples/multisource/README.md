# Multi Source

The values for configuration parameters are read from multiple sources.
The order they are applied is:

1. Default values are overwritten by
2. Parameters defined in the config-file, which are overwritten by
3. Environment variables, which are overwritten by
4. Command-Line parameters

## Naming Schema for Environment Variables

Each config parameter can be set via environment variable. The name of the according environment variable follows this schema.

### Schema

The following pseudo code lines out how to derive the name of an environment name from a given config parameter.

```go
// basic schema
nameOfEnvVar = <envPrefix>_<name of the config parameter>
// convert to upper case
nameOfEnvVar = toUpperCase(nameOfEnvVar)
// replace all occurrences of a dot and a dash by an underscore
nameOfEnvVar = replaceDotAndDashByUnderscore(nameOfEnvVar)
```

### envPrefix

Is the value in uppercase that was set when creating the provider instance.

In the following example `my-app` is the `envPrefix`

```go
prefixForEnvironmentVariables := "my-app"
provider, err := config.NewConfigProvider(
    &cfg,
    nameOfTheConfig,
    ...
)
```

### Name of the Config Parameter

The name of the config parameter is the complete name including all upper hierarchy levels.
Each config struct field is annotated with the cfg struct tag. With this annotation each field gets a name assigned. If there are fields of type struct that contain fields themselves these are concatenated by a dot (`.`).

In the following example

```go
type ThemeConfig struct {
    Header FormattedTextBox `cfg:"{'name':'header'}"`
}

type FormattedTextBox struct {
    Font   Font   `cfg:"{'name':'font'}"`
}

type Font struct {
    ForegroundColor string `cfg:"{'name':'foreground-color'}"`
}
```

The full name of the config parameter for the font color of the header is not just the name of the lowest level (`foreground-color`) but is the concatenation of all upper levels.
Which is `header.font.foreground-color`

### Example

```go
envPrefix := "my-app"
nameOfConfigParameter := "header.font.foreground-color"

// The according environment variable name is
nameOfEnvVar := MY_APP_HEADER_FONT_FOREGROUND_COLOR
```

## Config File Parameter

Per default the config file can be specified as command line argument `--config-file`.

But this behavior can be adjusted by using the option `config.CfgFile(..),`

```go
provider, err := config.NewConfigProvider(
    &cfg,
    nameOfTheConfig,
    prefixForEnvironmentVariables,
    config.CfgFile("config-file-name", "f"), // Overwrites the default parameter (--config-file) for the config file
)
```
