package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type HttpServer struct {
	Host string
	Port int
}

type Config struct {
	Env         string     `yaml:"env" env: "ENV" env-required: "true" env-default: "production"`
	StoragePath string     `yaml:"storage_path" env: "STORAGE_PATH" env-required: "true" env-default: "/storage"`
	HttpServer `yaml:"http_server"`
}

func MustConfig() *Config {
	var configPath string = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}
	if _, error := os.Stat(configPath); os.IsNotExist(error) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}
	var config Config
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	return &config
}
