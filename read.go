package config

import "fmt"

// ReadConfig parses commandline arguments, reads parameters from config and from environment
func (p *providerImpl) ReadConfig(args []string) error {

	if p.pFlagSet == nil {
		return fmt.Errorf("Pflag is nil")
	}

	if p.Viper == nil {
		return fmt.Errorf("Viper is nil")
	}

	if err := p.setDefaults(); err != nil {
		return err
	}

	if err := p.registerAndParseFlags(args); err != nil {
		return err
	}

	if err := p.registerEnvParams(); err != nil {
		return err
	}

	configFilenameEntryname := p.configFileEntry.name
	cfgFile := p.Viper.GetString(configFilenameEntryname)
	if err := p.readCfgFile(cfgFile); err != nil {
		return err
	}

	// apply the configuration automatically in case a target is available
	if p.configTarget != nil {
		return applyConfig(p, p.configTarget, "", configTag{}, p.mappingFuncRegistry)
	}

	return nil
}

func (p *providerImpl) readCfgFile(cfgFileName string) error {
	if len(cfgFileName) == 0 {
		return nil
	}
	p.Viper.SetConfigFile(cfgFileName)
	if err := p.Viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
