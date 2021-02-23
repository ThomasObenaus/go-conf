# MappingFunc

Sometimes it is necessary to intercept the parameter that is read from the configuration before it is filled into the config struct.
This can be useful for:

- handling complex data structures mapped from a string
- for validation purposes
- for mapping a string to another type

A mapping function has to be

1. defined
2. registered
3. referenced in the cfg tag of the according field in the config struct

## Define a Mapping Function

A mapping function is a function that maps a given input type into an output type.
**For structs all mapping functions are called with a string as input parameter.**

```go
// signature of a mapping function
type MappingFunc func(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error)
```

Since the mapping function deals with the interface{} type instead of real types it is important to know that:

- The type of the parameter `rawUntypedValue` is the type of the config parameter as it is read in and hence should be checked.
- The type of the returned value has to match the type of according the field in the config struct.

In the example below the given mapping function is intended to map a config parameter (int) to a struct field of type string.

```go
func mapIntToWorkingDay(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
    // ensure that the expectation of the input type is satisfied
    day, ok := rawUntypedValue.(int)
    if !ok {
        return nil, fmt.Errorf("Expected a int. Type '%T' is not supported", rawUntypedValue)
    }


    // return the target type (int8)
    switch day {
    case 0:
        return "Monday", nil
    case 1:
        return "Tuesday", nil
    case 2:
        return "Wednesday", nil
    case 3:
        return "Thursday", nil
    case 4:
        return "Friday", nil
    default:
        return "Unknown", fmt.Errorf("%d is not a working day (only 0-4 are supported days)",day)
    }
}
```

## Register the Mapping Function

Before reading the config one can register a mapping function at the provider. The name used for registration can then be used in step 3 to refer to that mapping function and actually link the struct field with that function.

```go
provider,_ := config.NewConfigProvider(...)

// Register the mapping function
_ = provider.RegisterMappingFunc("mapIntToWorkingDay", mapIntToWorkingDay)

_ = provider.ReadConfig(args)
```

## Link the Structfield and the Mapping Function

A registered mapping function can be linked to a config struct field. This can be done with the parameter `mapfun` in the cfg tag of that field it should be assigned to.

```go
type MyConfig struct {
    WorkDay string `cfg:"{'name':'workday','mapfun':'mapIntToWorkingDay','desc':'Workday as int.','default':0}"`
}
```
