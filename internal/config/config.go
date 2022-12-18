package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC struct {
		IP   string `yaml:"ip" env:"GRPC-IP"`
		Port int    `yaml:"port" env:"GRPC-PORT"`
	} `yaml:"grpc"`
	Writer struct {
		IP   string `yaml:"ip" env:"WRITER-IP"`
		Port string `yaml:"port" env:"WRITER-PORT"`
	} `yaml:"writer"`
	AppConfig struct {
		LogLevel string `yaml:"log-level" env:"LOG_LEVEL" env-default:"trace"`
	} `yaml:"app"`
	SQLite struct {
		DBPath string `yaml:"db_path" env:"DB_PATH" env-required:"true"`
	} `yaml:"sqlite"`
}

const (
	EnvConfigPathName  = "CONFIG-PATH"
	FlagConfigPathName = "config"
)

var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		flag.StringVar(&configPath, FlagConfigPathName, ".configs/config.local.yaml", "this is app config file")
		flag.Parse()

		log.Print("config init")

		if configPath == "" {
			configPath = os.Getenv(EnvConfigPathName)
		}

		if configPath == "" {
			log.Fatal("config path is required")
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			helpText := "Read Only"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
