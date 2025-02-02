package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	File string
	Type string
}

var EnvFile = &Config{
	File: ".env",
	Type: "env",
}

func ConfigFile(file string, type_ string) Config {
	return Config{
		File: file,
		Type: type_,
	}
}

func (c Config) LoadConfig(dst interface{}) error {
	cwd := os.Getenv("APP_ROOT")

	viper.AddConfigPath(cwd)
	viper.SetConfigName(c.File)
	viper.SetConfigType(c.Type)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(dst); err != nil {
		return err
	}

	return nil
}
