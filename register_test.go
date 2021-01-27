package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RegisterEnvParams(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("test", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	pImpl := toProviderImpl(t, provider)

	// WHEN
	err = pImpl.registerEnvParams()

	// THEN
	assert.NoError(t, err)
}

func Test_RegisterEnvParamsShouldFail(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	pImpl := toProviderImpl(t, provider)

	// WHEN
	err = pImpl.registerEnvParams()

	// THEN
	assert.Error(t, err)

	// GIVEN
	provider, err = NewProvider(entries, "configName", "envPrefix", CfgFile("", ""))
	require.NoError(t, err)
	pImpl = toProviderImpl(t, provider)

	// WHEN
	err = pImpl.registerEnvParams()

	// THEN
	assert.Error(t, err)
}

func Test_RegisterAndParseFlags(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	entries = append(entries, NewEntry("test2", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	args := []string{"--test1=A"}
	pImpl := toProviderImpl(t, provider)

	// WHEN
	err = pImpl.registerAndParseFlags(args)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "A", provider.GetString("test1"))
	assert.False(t, provider.IsSet("test2"))
}

func Test_RegisterAndParseFlags_ShouldFail(t *testing.T) {

	// GIVEN - unknown parameter
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	args := []string{"--unkown-param=A"}
	pImpl := toProviderImpl(t, provider)

	// WHEN
	err = pImpl.registerAndParseFlags(args)

	// THEN
	assert.Error(t, err)
	assert.False(t, provider.IsSet("test1"))

	// GIVEN - invalid entry
	entries = append(entries, NewEntry("", "usage"))
	provider, err = NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	pImpl = toProviderImpl(t, provider)
	args = []string{}

	// WHEN
	err = pImpl.registerAndParseFlags(args)

	// THEN
	assert.Error(t, err)
	assert.False(t, provider.IsSet("test1"))
}

func Test_SetDefaults(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage", Default("2h")))
	entries = append(entries, NewEntry("test2", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	pImpl := toProviderImpl(t, provider)

	// WHEN
	err = pImpl.setDefaults()

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, time.Hour*2, provider.GetDuration("test1"))
	assert.False(t, provider.IsSet("test2"))
}
