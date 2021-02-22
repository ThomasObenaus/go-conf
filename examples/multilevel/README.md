# Multi Level

go-conf supports multiple levels/ hierarchical configuration parameters. To do so one can just define a config structure that contains fields that are no primitive types but also structures.

## Example

The goal is to set the value of the headers foreground color.

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

## Hierarchy Separator

### Config File

The hierarchy is represented using the structural elements provided by yaml.

To set the value for the headers foreground color one has to specify:

```yaml
header:
  font:
    foreground-color: #FF0000
```

### Environment Variable

The hierarchy in environment variables is represented by an underscore. But the underscore is also used as replacement for the dash. Hence the hierarchy can't be safely deduced from the name of the environment variable.

To set the value for the headers foreground color one has to specify:

```bash
# Assumption: MY_APP was the envPrefix for the provider
export MY_APP_HEADER_FONT_FOREGROUND_COLOR=#FF0000
```

### Command Line Argument

The hierarchy for command line arguments is represented by a dot.

To set the value for the headers foreground color one has to specify:

```bash
go run . --header.font.foreground-color=#FF0000
```
