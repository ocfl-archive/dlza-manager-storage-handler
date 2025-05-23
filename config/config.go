package config

import (
	"emperror.dev/errors"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/je4/filesystem/v3/pkg/vfsrw"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/stashconfig"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"io/fs"
	"maps"
	"os"
)

type Config struct {
	ErrorConfig             string             `toml:"errorconfig"`
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
	TusServer               TusServer          `toml:"tusserver"`
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

type TusServer struct {
	Addr    string   `toml:"addr"`
	ExtAddr string   `toml:"extaddr"`
	TLSCert string   `toml:"tlscert"`
	TLSKey  string   `toml:"tlskey"`
	RootCA  []string `toml:"rootca"`
}

type Connection struct {
	Folder string
	VFS    vfsrw.Config
}

func LoadConfig(fSys fs.FS, fp string, conf *Config) error {
	if _, err := fs.Stat(fSys, fp); err != nil {
		path, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "cannot get current working directory")
		}
		fSys = os.DirFS(path)
		fp = "storagehandler.toml"
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
	err = os.Setenv("AWS_REQUEST_CHECKSUM_CALCULATION", "when_required")
	if err != nil {
		return errors.Wrapf(err, "cannot set env variable file while LoadConfig")
	}
	return nil
}

func LoadVfsConfig(storageLocations *pb.StorageLocations, cfg Config) (vfsrw.Config, error) {
	vfsMap := make(map[string]*vfsrw.VFS)
	for _, storageLocation := range storageLocations.StorageLocations {
		connection := Connection{}
		err := json.Unmarshal([]byte(storageLocation.Connection), &connection)
		if err != nil {
			return nil, errors.Wrapf(err, "error mapping json for storage location connection field")
		}
		maps.Copy(vfsMap, connection.VFS)
	}
	maps.Copy(vfsMap, getVfsTempMap(cfg))
	return vfsMap, nil
}

func getVfsTempMap(cfg Config) map[string]*vfsrw.VFS {
	vfsTemp := vfsrw.VFS{
		Type: cfg.S3TempStorage.Type,
		Name: cfg.S3TempStorage.Name,
		S3: &vfsrw.S3{
			AccessKeyID:     config.EnvString(cfg.S3TempStorage.Key),
			SecretAccessKey: config.EnvString(cfg.S3TempStorage.Secret),
			Endpoint:        config.EnvString(cfg.S3TempStorage.Url),
			Region:          "us-east-1",
			UseSSL:          true,
			Debug:           cfg.S3TempStorage.Debug,
			CAPEM:           cfg.S3TempStorage.CAPEM,
		},
	}

	tempVfsMap := make(map[string]*vfsrw.VFS)
	tempVfsMap[cfg.S3TempStorage.Name] = &vfsTemp
	return tempVfsMap
}
