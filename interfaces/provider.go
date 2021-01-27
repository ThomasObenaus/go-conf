package interfaces

import (
	"reflect"
	"time"
)

type Provider interface {
	/*
		RegisterMappingFunc used to register a function that will map the value and type provided as command line parameter or default value.
		This handles the case where the type of a field defined in the config annotation does not match the type of the field that is annotated.

		Example:
		 type cfg1 struct {
		  LogLevel zerolog.Level `cfg:"{'name':'logl','default':'info','mapfun':'strToLogLevel'}"`
		 }
		 // [..]
		 provider.RegisterMappingFunc("toUpperCase",func(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error){
			 asStr, ok := rawUntypedValue.(string)
			 if !ok {
			 	return nil, fmt.Errorf("Expected a string. Type '%T' is not supported", rawUntypedValue)
			 }
			 return zerolog.ParseLevel(asStr)
		 })
		Here F1 is of type zerolog.Level (int8) and the defined type in the annotation is string (based on the default value).
		In order to support this situation we have to apply the defined mapping functions.

		A mapping function can be used also to do some conversions. For example if string value just should be converted to upper case.

		Example:
		 type cfg2 struct {
		  LogLevel string `cfg:"{'name':'logl','default':'info','mapfun':'toUpperCase'}"`
		 }
		 // [..]
		 provider.RegisterMappingFunc("toUpperCase",func(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error){
			 asStr,ok := rawUntypedValue.(string)
			 if !ok {
				 return nil, fmt.Errorf("Value must be a string")
			 }
			 return strings.ToUpper(asStr),nil
		 })
	*/
	RegisterMappingFunc(name string, mFunc MappingFunc) error

	ReadConfig(args []string) error

	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetSizeInBytes(key string) uint
	IsSet(key string) bool
	String() string
	Log(lvl LogLevel, formatString string, a ...interface{})
	Usage() string
}

/*
	MappingFunc type that specifies a mapping function.

	rawUntypedValue - the incoming value

	targetType - the type the returned result should have
*/
type MappingFunc func(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error)
