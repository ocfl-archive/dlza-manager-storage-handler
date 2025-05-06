package models

import (
	"emperror.dev/errors"
	"encoding/json"
	"github.com/je4/filesystem/v3/pkg/vfsrw"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"maps"
)

type StorageHandler struct {
	LogLevel string
	LogFile  string
	Addr     string
	VFS      vfsrw.Config
}

func LoadStorageHandlerConfig(storageLocations []*pb.StorageLocation) (*StorageHandler, error) {
	vfsMap := make(map[string]*vfsrw.VFS)
	for _, storageLocation := range storageLocations {
		connection := Connection{}
		err := json.Unmarshal([]byte(storageLocation.Connection), &connection)
		if err != nil {
			return nil, errors.Wrapf(err, "error mapping json for storageLocation: %v", storageLocation.Alias)
		}
		maps.Copy(vfsMap, connection.VFS)
	}

	var config = &StorageHandler{LogLevel: "DEBUG", Addr: "localhost:8080", VFS: vfsMap}

	return config, nil
}
