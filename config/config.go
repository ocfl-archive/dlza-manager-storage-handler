package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/stashconfig"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"io/fs"
	"os"
)

type Config struct {
	LocalAddr               string             `toml:"localaddr"`
	Domain                  string             `toml:"domain"`
	ExternalAddr            string             `toml:"externaladdr"`
	Bearer                  string             `toml:"bearer"`
	ResolverAddr            string             `toml:"resolveraddr"`
	ResolverTimeout         config.Duration    `toml:"resolvertimeout"`
	ResolverNotFoundTimeout config.Duration    `toml:"resolvernotfoundtimeout"`
	ActionTimeout           config.Duration    `toml:"actiontimeout"`
	ServerTLS               *loader.Config     `toml:"server"`
	ClientTLS               *loader.Config     `toml:"client"`
	GRPCClient              map[string]string  `toml:"grpcclient"`
	Log                     stashconfig.Config `toml:"log"`
	S3TempStorage           S3TempStorage      `toml:"s3tempstorage"`
}

type S3TempStorage struct {
	Type         string `yaml:"type" toml:"type"`
	Name         string `yaml:"name" toml:"name"`
	Key          string `yaml:"key" toml:"key"`
	Secret       string `yaml:"secret" toml:"secret"`
	Bucket       string `yaml:"bucket" toml:"bucket"`
	ApiUrlValue  string `yaml:"api-url-value" toml:"apiurlvalue"`
	UploadFolder string `yaml:"upload-folder" toml:"uploadfolder"`
	Url          string `yaml:"url" toml:"url"`
	CAPEM        string `yaml:"capem" toml:"capem"`
	Debug        bool   `yaml:"debug" toml:"debug"`
}

func LoadConfig(fSys fs.FS, fp string, conf *Config) error {
	if _, err := fs.Stat(fSys, fp); err != nil {
		path, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "cannot get current working directory")
		}
		fSys = os.DirFS(path)
		fp = "mediaserveraction.toml"
	}
	data, err := fs.ReadFile(fSys, fp)
	if err != nil {
		return errors.Wrapf(err, "cannot read file [%v] %s", fSys, fp)
	}
	_, err = toml.Decode(string(data), conf)
	if err != nil {
		return errors.Wrapf(err, "error loading config file %v", fp)
	}
	if conf.S3TempStorage.Secret == "" {
		conf.S3TempStorage.Secret = os.Getenv("S3_SECRET")
	}
	return nil
}
