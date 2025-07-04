package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

type Config struct {
	App struct {
		WebhookURL       string `yaml:"webhook_url" mapstructure:"webhook_url"`
		WebhookKey       string `yaml:"webhook_key" mapstructure:"webhook_key"` //optional
		SendIntervalSecs int    `yaml:"send_interval_seconds" mapstructure:"send_interval_seconds"`
		MessageCharLimit int    `yaml:"message_char_limit" mapstructure:"message_char_limit"`
		MaxRetries       int    `yaml:"max_retries" mapstructure:"max_retries"`
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
		URL string `yaml:"url" mapstructure:"url"` //optional
	} `yaml:"redis" mapstructure:"redis"`
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	var config Config

	// Load from YAML
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config read error: %v", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Config unmarshal error: %v", err)
	}

	// Overwrite from ENV if available
	overrideString(&config.App.WebhookURL, "WEBHOOK_URL")
	overrideString(&config.App.WebhookKey, "WEBHOOK_KEY")

	overrideString(&config.Database.Host, "PGHOST")
	overrideInt(&config.Database.Port, "PGPORT")
	overrideString(&config.Database.User, "PGUSER")
	overrideString(&config.Database.Password, "PGPASSWORD")
	overrideString(&config.Database.Name, "PGDATABASE")
	overrideString(&config.Database.SSLMode, "PGSSLMODE")

	overrideString(&config.Redis.URL, "REDIS_URL")

	return config
}

func overrideString(field *string, envKey string) {
	if val := os.Getenv(envKey); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envKey string) {
	if val := os.Getenv(envKey); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			*field = intVal
		}
	}
}
