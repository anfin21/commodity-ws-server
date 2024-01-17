package configs

func DefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Name:       DefaultServerName,
			Env:        DefaultServerEnv,
			Host:       DefaultServerHost,
			Port:       DefaultServerPort,
			ConfigFile: DefaultServerConfigFile,
			Timezone:   DefaultServerTimezone,
		},
	}
}

const (
	DefaultServerName                 = ""
	DefaultServerEnv        ServerEnv = ServerEnvDevelopment
	DefaultServerHost                 = ""
	DefaultServerPort                 = 8089
	DefaultServerConfigFile           = "configs/dev.config.yaml"
	DefaultServerTimezone             = "Asia/Ho_Chi_Minh"
)
