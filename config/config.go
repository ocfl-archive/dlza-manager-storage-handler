package config

import (
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	"log"
	"os"

	"github.com/jinzhu/configor"
)

type Service struct {
	ServiceName string `yaml:"service_name" toml:"ServiceName"`
	Host        string `yaml:"host" toml:"Host"`
	Port        int    `yaml:"port" toml:"Port"`
}

type Logging struct {
	LogLevel string
	LogFile  string
}

type Config struct {
	ServerConfig   models.ServerConfig `yaml:"server-config" toml:"ServerConfig"`
	Handler        Service             `yaml:"handler" toml:"Handler"`
	StorageHandler Service             `yaml:"storage-handler" toml:"StorageHandler"`
	Clerk          Service             `yaml:"clerk" toml:"Clerk"`
	S3TempStorage  S3TempStorage       `yaml:"s3-temp-storage" toml:"S3TempStorage"`
	Logging        Logging             `yaml:"logging" toml:"Logging"`
}

type S3TempStorage struct {
	Type         string `yaml:"type" toml:"Type"`
	Name         string `yaml:"name" toml:"Name"`
	Key          string `yaml:"key" toml:"Key"`
	Secret       string `yaml:"secret" toml:"Secret"`
	Bucket       string `yaml:"bucket" toml:"Bucket"`
	ApiUrlValue  string `yaml:"api-url-value" toml:"ApiUrlValue"`
	UploadFolder string `yaml:"upload-folder" toml:"UploadFolder"`
	Url          string `yaml:"url" toml:"Url"`
	CAPEM        string `yaml:"capem" toml:"CAPEM"`
	Debug        bool   `yaml:"debug" toml:"Debug"`
}

// GetConfig creates a new config from a given environment
func GetConfig(configFile string) Config {
	conf := Config{}
	if configFile == "" {
		configFile = "config.yml"
	}
	err := configor.Load(&conf, configFile)
	if err != nil {
		log.Fatal(err)
	}
	if conf.S3TempStorage.Secret == "" {
		conf.S3TempStorage.Secret = os.Getenv("S3_SECRET")
	}
	return conf
}
