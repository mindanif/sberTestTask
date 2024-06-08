package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.BindEnv("database.url", "DATABASE_URL")
	viper.BindEnv("server.port", "SERVER_PORT")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
