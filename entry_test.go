package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewEntry_SimpleNoDefault(t *testing.T) {

	entry := NewEntry("help", "Print some help")
	assert.True(t, entry.bindEnv)
	assert.True(t, entry.bindFlag)
	assert.Equal(t, "help", entry.name)
	assert.Equal(t, "Print some help", entry.usage)
	assert.Nil(t, entry.defaultValue)
	assert.Empty(t, entry.flagShortName)
}

func Test_NewEntry_SimpleWithDefault(t *testing.T) {

	entry := NewEntry("port", "the port", Default(8080))
	assert.True(t, entry.bindEnv)
	assert.True(t, entry.bindFlag)
	assert.Equal(t, "port", entry.name)
	assert.Equal(t, "the port", entry.usage)
	assert.Equal(t, 8080, entry.defaultValue)
	assert.Empty(t, entry.flagShortName)
}

func Test_NewEntry_Full(t *testing.T) {

	entry := NewEntry("port", "the port", Default(8080), ShortName("p"), Bind(false, false))
	assert.False(t, entry.bindEnv)
	assert.False(t, entry.bindFlag)
	assert.Equal(t, "port", entry.name)
	assert.Equal(t, "the port", entry.usage)
	assert.Equal(t, "p", entry.flagShortName)
	assert.Equal(t, 8080, entry.defaultValue)
}

func Test_CheckViper(t *testing.T) {

	err := checkViper(nil, Entry{})
	assert.Error(t, err)

	vp := viper.New()
	require.NotNil(t, vp)

	err = checkViper(vp, Entry{})
	assert.Error(t, err)

	cfgE := Entry{
		name: "bla",
	}
	err = checkViper(vp, cfgE)
	assert.NoError(t, err)
}

func Test_SetDefault_OK(t *testing.T) {

	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := Entry{
		name:         "bla",
		defaultValue: 20,
	}
	err := setDefault(vp, cfgE)
	assert.NoError(t, err)

	assert.NotNil(t, vp.GetInt(cfgE.name))
	assert.Equal(t, cfgE.defaultValue, vp.GetInt(cfgE.name))
}

func Test_RegisterEnv_OK(t *testing.T) {

	envPrefix := "ABCD"
	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := Entry{
		name:        "flag",
		bindEnv:     true,
		desiredType: reflect.TypeOf(""),
	}
	err := registerEnv(vp, envPrefix, cfgE)
	assert.NoError(t, err)
	os.Setenv(envPrefix+"_"+strings.ToUpper(cfgE.name), "test1")
	assert.NotEmpty(t, vp.Get(cfgE.name))
	assert.Equal(t, "test1", vp.Get(cfgE.name))

	cfgE = Entry{
		name:        "flag",
		bindEnv:     true,
		desiredType: reflect.TypeOf(""),
	}
	err = registerEnv(vp, envPrefix, cfgE)
	assert.NoError(t, err)
	os.Setenv(strings.ToUpper(envPrefix+"_"+cfgE.name), "test2")
	assert.NotEmpty(t, vp.Get(cfgE.name))
	assert.Equal(t, "test2", vp.Get(cfgE.name))

	// without env prefix
	vp = viper.New()
	require.NotNil(t, vp)
	err = registerEnv(vp, "", cfgE)
	assert.NoError(t, err)
	os.Setenv(strings.ToUpper(cfgE.name), "test3")
	assert.Equal(t, "test3", vp.Get(cfgE.name))
}

func Test_RegisterEnv_Fail(t *testing.T) {

	err := registerEnv(nil, "ABC", Entry{})
	assert.NoError(t, err)

	err = registerEnv(nil, "ABC", Entry{bindEnv: true})
	assert.Error(t, err)

	vp := viper.New()
	require.NotNil(t, vp)

	cfgE := Entry{bindEnv: true}
	err = registerEnv(vp, "ABC", cfgE)
	assert.Error(t, err)
}

func Test_RegisterFlag_Fail(t *testing.T) {

	err := registerFlag(nil, Entry{})
	assert.NoError(t, err)

	err = registerFlag(nil, Entry{bindFlag: true})
	assert.Error(t, err)

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	require.NotNil(t, flagSet)

	cfgE := Entry{bindFlag: true}
	err = registerFlag(flagSet, cfgE)
	assert.Error(t, err)

}

func Test_RegisterFlag_Ok(t *testing.T) {

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	require.NotNil(t, flagSet)

	// String
	cfgE := Entry{
		bindFlag:      true,
		name:          "flag1",
		defaultValue:  "default",
		usage:         "The default value",
		flagShortName: "a",
		desiredType:   reflect.TypeOf(""),
	}
	err := registerFlag(flagSet, cfgE)
	require.NoError(t, err)
	flag := flagSet.Lookup(cfgE.name)
	require.NotNil(t, flag)
	assert.Equal(t, cfgE.defaultValue.(string), flag.DefValue)
	assert.Equal(t, cfgE.flagShortName, flag.Shorthand)
	assert.Equal(t, cfgE.usage, flag.Usage)

	// Uint
	cfgE = Entry{
		bindFlag:      true,
		name:          "flag2",
		defaultValue:  uint(1),
		usage:         "An uint",
		flagShortName: "b",
	}
	err = registerFlag(flagSet, cfgE)
	assert.NoError(t, err)
	flag = flagSet.Lookup(cfgE.name)
	require.NotNil(t, flag)
	assert.Equal(t, fmt.Sprintf("%d", cfgE.defaultValue.(uint)), flag.DefValue)
	assert.Equal(t, cfgE.flagShortName, flag.Shorthand)
	assert.Equal(t, cfgE.usage, flag.Usage)

	// Int
	cfgE = Entry{
		bindFlag:      true,
		name:          "flag3",
		defaultValue:  int(1),
		usage:         "An int",
		flagShortName: "c",
	}
	err = registerFlag(flagSet, cfgE)
	assert.NoError(t, err)
	flag = flagSet.Lookup(cfgE.name)
	require.NotNil(t, flag)
	assert.Equal(t, fmt.Sprintf("%d", cfgE.defaultValue.(int)), flag.DefValue)
	assert.Equal(t, cfgE.flagShortName, flag.Shorthand)
	assert.Equal(t, cfgE.usage, flag.Usage)

	// Type not supported
	typeNotSupported := reflect.TypeOf("")

	cfgE = Entry{
		bindFlag:      true,
		name:          "flag4",
		defaultValue:  typeNotSupported,
		usage:         "Reflect type info",
		flagShortName: "d",
	}
	err = registerFlag(flagSet, cfgE)
	assert.Error(t, err)
}

func ExampleNewEntry() {
	entry := NewEntry("port", "The port of the service", Default(8080), ShortName("p"))
	fmt.Printf("%s", entry)
	// Output:
	// --port (-p) [default:8080 (int)]	- The port of the service
}
