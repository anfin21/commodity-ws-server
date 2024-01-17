package configs

type ServerEnv string

const (
	ServerEnvLocalhost   ServerEnv = "local"
	ServerEnvDevelopment ServerEnv = "dev"
	ServerEnvStaging     ServerEnv = "staging"
	ServerEnvProduction  ServerEnv = "prod"
)

type Config struct {
	Server *ServerConfig `mapstructure:"SERVER"`
}

type ServerConfig struct {
	Name       string    `mapstructure:"NAME"`
	Env        ServerEnv `mapstructure:"ENV"`
	ConfigFile string    `mapstructure:"CONFIG_FILE"`
	Host       string    `mapstructure:"HOST"`
	Port       int       `mapstructure:"PORT"`
	Timezone   string    `mapstructure:"TZ"`
}

func (cfg *Config) GetServerEnv() ServerEnv {
	if cfg.Server.Env == "" {
		return DefaultServerEnv
	}
	return cfg.Server.Env
}
