package module

import (
	"github.com/spf13/viper"
)

type Configure struct {
	*viper.Viper
	cfgFile *string
}

func (p *Configure) From(v *viper.Viper) *Configure {
	p.Viper = v
	return p
}

func (p *Configure) CfgFile(cfgFile *string) *Configure {
	p.cfgFile = cfgFile
	return p
}

func (p *Configure) ReadInConfig() error {
	if p.cfgFile != nil && len(*p.cfgFile) > 0 {
		p.SetConfigFile(*p.cfgFile)
	} else {
		p.SetConfigName("config")
		p.AddConfigPath(".")
		p.AddConfigPath("$HOME")
		p.SetConfigType("yaml")
	}
	if err := p.Viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			return nil
		default:
			return err
		}
	}
	return nil
}
