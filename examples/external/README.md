# External Types

External types are types that can't be modified since they **are defined in an external library**. Hence **their fields can't be annotated with the cfg tag**.

To be able to use them in the config and automatically populate the config parameter value to the config struct, two things have to be done:

1. The field has to be annotated with the cfg tag.
2. A mapping function mapping from a string to the field type has to be defined and registered.

**For structs all mapping functions are called with a string as input parameter.**

## Annotate the Field

The field `Field1` uses the external type `time.Time`. It is annotated with the with the cfg tag.
With the tag definition the mapping function strToTime is referenced.

```go
type ExternalTypes struct {
    Field1 time.Time `cfg:"{'name':'field-1','mapfun':'strToTime'}"`
}
```

## Define and Register a Mapping Function

As an example the following function maps a given string to a `time.Time`

```go
// strToTime maps a string to a time.Time
func strToTime(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
    asStr, ok := rawUntypedValue.(string)
    if !ok {
        return nil, fmt.Errorf("Expected a string. Type '%T' is not supported", rawUntypedValue)
    }
    t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", asStr)
    if err != nil {
        return nil, fmt.Errorf("Parse %s to time failed: %s", asStr, err.Error())
    }
    return t, nil
}
```

As a last step that function has to be registered at the provider.

```go
[...]
if err := provider.RegisterMappingFunc("strToTime", strToTime); err != nil {
    panic(err)
}
[...]
```
