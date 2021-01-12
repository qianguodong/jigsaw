package settings

type Option interface {
	Apply(settings *Settings)
}

func WithConfigFile(configFile string) Option {
	return withConfigFile{configFile: configFile}
}

type withConfigFile struct {
	configFile string
}

func (o withConfigFile) Apply(settings *Settings) {
	settings.configFile = o.configFile
}
