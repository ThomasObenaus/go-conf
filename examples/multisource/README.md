# Multi Source

The values for configuration parameters are read from multiple sources.
The order they are applied is:

1. Default values are overwritten by
2. Parameters defined in the config-file, which are overwritten by
3. Environment variables, which are overwritten by
4. Command-Line parameters

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
