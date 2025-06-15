package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App struct {
		WebhookURL       string `yaml:"webhook_url" mapstructure:"webhook_url"`
		WebhookKey       string `yaml:"webhook_key" mapstructure:"webhook_key"` //in case you have
		SendIntervalSecs int    `yaml:"send_interval_seconds" mapstructure:"send_interval_seconds"`
	} `yaml:"app" mapstructure:"app"`

	Database struct {
		Host     string `yaml:"host" mapstructure:"host"`
		Port     int    `yaml:"port" mapstructure:"port"`
		User     string `yaml:"user" mapstructure:"user"`
		Password string `yaml:"password" mapstructure:"password"`
		Name     string `yaml:"name" mapstructure:"name"`
		SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
	} `yaml:"database" mapstructure:"database"`

	Redis struct {
		Addr string `yaml:"addr" mapstructure:"addr"`
		DB   int    `yaml:"db" mapstructure:"db"`
	} `yaml:"redis" mapstructure:"redis"`
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config read error: %v", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Config unmarshal error: %v", err)
	}
	return config
}
