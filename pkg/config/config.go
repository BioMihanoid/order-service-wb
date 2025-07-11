package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DbConfig DbConfig     `mapstructure:"db"`
	Server   ServerConfig `mapstructure:"server"`
	Cache    CacheConfig  `mapstructure:"cache"`
	Kafka    KafkaConfig  `mapstructure:"kafka"`
}

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type CacheConfig struct {
	Size int `mapstructure:"size"`
}

type KafkaConfig struct {
	Broker string `mapstructure:"broker"`
	Topic  string `mapstructure:"topic"`
	Group  string `mapstructure:"group"`
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	configName := os.Getenv("CONFIG_NAME")
	if configPath == "" || configName == "" {
		log.Fatal("CONFIG_PATH or CONFIG_NAME environment variable not set")
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error unmarshalling config")
	}

	return &config
}

func (d *DbConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.Username, d.Password, d.Database, d.SSLMode)
}
