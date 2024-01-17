package configs

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	config := DefaultConfig()
	setDefault(config)

	if err := config.load(); err != nil {
		return config, err
	}

	return config, nil
}

func setDefault(defaultConfig *Config) {
	viper.SetDefault("SERVER", defaultConfig.Server)
}

func (cfg *Config) load() error {
	viper.SetConfigFile(configFilePath())

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.AutomaticEnv()
	if err := viper.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.ZeroFields = false
	}); err != nil {
		return err
	}
	return nil
}

func configFilePath() string {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return DefaultServerConfigFile
	}
	return configFile
}
