package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReadCfgFile(t *testing.T) {

	// GIVEN
	configFilename := "test/data/config.yaml"
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	entries = append(entries, NewEntry("test2", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)

	// WHEN
	pImpl := toProviderImpl(t, provider)
	err = pImpl.readCfgFile(configFilename)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "A", provider.GetString("test1"))
	assert.False(t, provider.IsSet("test2"))
	assert.Equal(t, configFilename, pImpl.Viper.ConfigFileUsed())
}

func Test_ReadCfgFile_AllowNoCfgFile(t *testing.T) {

	// GIVEN
	configFilename := ""
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)

	// WHEN
	pImpl := toProviderImpl(t, provider)
	err = pImpl.readCfgFile(configFilename)

	// THEN
	assert.NoError(t, err)
	assert.False(t, provider.IsSet("test1"))
	assert.Empty(t, pImpl.Viper.ConfigFileUsed())
}

func Test_ReadCfgFile_ShouldFail(t *testing.T) {

	// GIVEN
	configFilename := "does_not_exist.yaml"
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)

	// WHEN
	pImpl := toProviderImpl(t, provider)
	err = pImpl.readCfgFile(configFilename)

	// THEN
	assert.Error(t, err)
	assert.False(t, provider.IsSet("test1"))
	assert.Equal(t, configFilename, pImpl.Viper.ConfigFileUsed())
}

func Test_ReadConfig_ShouldFail(t *testing.T) {

	// GIVEN
	var entries []Entry
	provider, err := NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	pImpl := toProviderImpl(t, provider)
	args := []string{}
	pImpl.Viper = nil

	// WHEN
	err = provider.ReadConfig(args)

	// THEN
	assert.Error(t, err)

	// GIVEN
	provider, err = NewProvider(entries, "configName", "envPrefix")
	require.NoError(t, err)
	pImpl = toProviderImpl(t, provider)
	pImpl.pFlagSet = nil

	// WHEN
	err = provider.ReadConfig(args)

	// THEN
	assert.Error(t, err)
}
