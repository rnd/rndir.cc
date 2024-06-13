package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	configPath = "env/config"
	configType = "ini"
)

type Config struct {
	Server   ServerCfg
	Postgres PostgresCfg
}

type ServerCfg struct {
	Name        string
	Version     string
	Environment string

	Host string
	Port string
}

type PostgresCfg struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load will loads default config from file, replace matching key from
// environment variables and unmarshal it to config struct.
func Load() (Config, error) {
	cfg := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	cfg.SetConfigFile(configPath)
	cfg.SetConfigType(configType)
	cfg.AutomaticEnv()

	var c Config
	if err := cfg.ReadInConfig(); err != nil {
		return c, fmt.Errorf("failed to read environment config")
	}
	if err := cfg.Unmarshal(&c); err != nil {
		return c, fmt.Errorf("failed to unmarshal environment config")
	}
	return c, nil
}
