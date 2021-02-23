# Required

All fields that are annotated with the cfg tag but have no default value specified are treated as required variables.
If for those parameters no value is given, neither via config-file nor via environment variable or command line argument, then reading the config will fail with an error.

```go
type MyFontConfig struct {
    Color string `cfg:"{'name':'color','desc':'A required parameter (no default is defined)'}"`
}
```
