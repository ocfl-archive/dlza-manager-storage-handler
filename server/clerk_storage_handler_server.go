package server

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
	"github.com/je4/filesystem/v2/pkg/writefs"
	lm "github.com/je4/utils/v2/pkg/logger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
)

type ClerkStorageHandlerServer struct {
	storageHandlerPb.UnimplementedClerkStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
}

const LOGFORMAT = `%{time:2006-01-02T15:04:05.000} %{shortpkg}::%{longfunc} [%{shortfile}] > %{level:.5s} - %{message}`

func (c *ClerkStorageHandlerServer) CreateStoragePartition(ctx context.Context, storagePartition *pb.StoragePartition) (*pb.Status, error) {

	storageLocation, err := c.ClientStorageHandlerHandler.GetStorageLocationById(ctx, &pb.Id{Id: storagePartition.StorageLocationId})
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "Could not get storage location by ID: %v", storagePartition.StorageLocationId)
	}
	storagePartitionWithAlias, err := c.ClientStorageHandlerHandler.GetAndSaveStoragePartitionWithRelevantAlias(ctx, storagePartition)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "Could not get and save storage partition with alias for storage location with ID: %v", storagePartition.StorageLocationId)
	}
	storageLocations := make([]*pb.StorageLocation, 0)
	config, err := models.LoadStorageHandlerConfig(append(storageLocations, storageLocation))
	daLogger, lf := lm.CreateLogger("storage-handler", string(config.LogFile), nil, string(config.LogLevel), LOGFORMAT)
	defer lf.Close()

	connection := models.Connection{}
	err = json.Unmarshal([]byte(storageLocation.Connection), &connection)
	if err != nil {
		return nil, errors.Wrapf(err, "error mapping json for storageLocation: %v", storageLocation.Alias)
	}

	vfs, err := vfsrw.NewFS(config.VFS, daLogger)
	if err != nil {
		daLogger.Errorf("cannot create vfs: %v", err)
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot create vfs: %v", err)
	}

	if err := writefs.MkDir(vfs, connection.Folder+storagePartitionWithAlias.Alias); err != nil {
		daLogger.Errorf("error creating partition")
		return &pb.Status{Ok: false}, errors.Wrapf(err, "error creating partition with alias: %v", storagePartition.Alias)
	}
	return &pb.Status{Ok: true}, nil
}
