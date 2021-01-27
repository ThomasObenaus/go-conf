package config

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func (p *providerImpl) registerEnvParams() error {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	p.Viper.SetEnvKeyReplacer(replacer)

	// register also the config file parameter
	if err := registerEnv(p.Viper, p.envPrefix, p.configFileEntry); err != nil {
		return err
	}

	for _, entry := range p.configEntries {
		if err := registerEnv(p.Viper, p.envPrefix, entry); err != nil {
			return err
		}
	}

	return nil
}

func (p *providerImpl) registerAndParseFlags(args []string) error {

	// register also the config file parameter
	if err := registerFlag(p.pFlagSet, p.configFileEntry); err != nil {
		return err
	}

	for _, entry := range p.configEntries {
		if err := registerFlag(p.pFlagSet, entry); err != nil {
			return err
		}
	}

	if err := p.pFlagSet.Parse(args); err != nil {

		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		return err
	}
	return p.Viper.BindPFlags(p.pFlagSet)
}

func (p *providerImpl) setDefaults() error {

	// regard also the config file parameter
	if err := setDefault(p.Viper, p.configFileEntry); err != nil {
		return err
	}

	for _, entry := range p.configEntries {
		if err := setDefault(p.Viper, entry); err != nil {
			return err
		}
	}
	return nil
}
