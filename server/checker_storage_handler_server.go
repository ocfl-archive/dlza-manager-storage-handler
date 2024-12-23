package server

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
	"github.com/je4/filesystem/v2/pkg/writefs"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/dlza-manager/models"
	"io"
	"path/filepath"
)

type CheckerStorageHandlerServer struct {
	storageHandlerPb.UnimplementedCheckerStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
	Logger                      zLogger.ZLogger
}

func (c *CheckerStorageHandlerServer) GetObjectInstanceChecksum(ctx context.Context, objectInstance *pb.ObjectInstance) (*pb.Id, error) {
	storageLocation, err := c.ClientStorageHandlerHandler.GetStorageLocationByObjectInstanceId(ctx, &pb.Id{Id: objectInstance.Id})
	if err != nil {
		c.Logger.Error().Msgf("cannot get storage location for object instance id %v, %s", objectInstance.Id, err)
		return nil, errors.Wrapf(err, "cannot get storage location for object instance id %v", objectInstance.Id)
	}
	vfsConfig, err := config.LoadVfsConfig(storageLocation.Connection)
	if err != nil {
		c.Logger.Error().Msgf("error mapping json for storage location connection field: %v", err)
		return nil, errors.Wrapf(err, "error mapping json for storage location connection field")
	}
	daLogger := zLogger.NewZWrapper(c.Logger)
	vfs, err := vfsrw.NewFS(vfsConfig, daLogger)

	sourceFP, err := vfs.Open(objectInstance.Path)
	if err != nil {
		c.Logger.Error().Msgf("cannot read file '%s': %v", objectInstance.Path, err)
		return nil, errors.Wrapf(err, "cannot read file '%v'", objectInstance.Path)
	}

	targetFP := io.Discard
	csWriter, err := checksum.NewChecksumWriter(
		[]checksum.DigestAlgorithm{checksum.DigestSHA512},
		targetFP,
	)
	if err != nil {
		c.Logger.Error().Msgf("cannot create new checksum writer: '%s'", err)
		return nil, errors.Wrapf(err, "cannot create new checksum writer: '%s'", err)
	}
	_, err = io.Copy(csWriter, sourceFP)
	if err != nil {
		c.Logger.Error().Msgf("error writing file")
		if err := csWriter.Close(); err != nil {
			c.Logger.Error().Msgf("cannot close checksum writer: %v", err)
		}
		if err := sourceFP.Close(); err != nil {
			c.Logger.Error().Msgf("cannot close source: %v", err)
		}
		return nil, errors.Wrapf(err, "error writing file")
	}
	if err := csWriter.Close(); err != nil {
		c.Logger.Error().Msgf("cannot close checksum writer: %v", err)
	}
	checksums, err := csWriter.GetChecksums()
	if err != nil {
		c.Logger.Error().Msgf("cannot get checksum for file '%v': %s", objectInstance.Path, err)
		if err := sourceFP.Close(); err != nil {
			c.Logger.Error().Msgf("cannot close source: %v", err)
		}
		c.Logger.Error().Msgf("cannot get checksum: %v", err)
		return nil, errors.Wrapf(err, "cannot get checksum for file '%v'", objectInstance.Path)
	}

	return &pb.Id{Id: checksums[checksum.DigestSHA512]}, nil
}

func (c *CheckerStorageHandlerServer) CopyArchiveTo(ctx context.Context, copyFromTo *pb.CopyFromTo) (*pb.NoParam, error) {
	storagePartition, err := c.ClientStorageHandlerHandler.GetStoragePartitionForLocation(ctx, &pb.SizeAndId{Size: copyFromTo.ObjectInstance.Size, Id: copyFromTo.LocationCopyTo.Id})
	if err != nil {
		c.Logger.Error().Msgf("cannot get storagePartition for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
		return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	storageLocationToCopyFrom, err := c.ClientStorageHandlerHandler.GetStorageLocationByObjectInstanceId(ctx, &pb.Id{Id: copyFromTo.ObjectInstance.Id})
	if err != nil {
		c.Logger.Error().Msgf("cannot get GetStorageLocationByObjectInstanceId for object instance ID: %v", copyFromTo.ObjectInstance.Id)
		return &pb.NoParam{}, errors.Wrapf(err, "cannot get GetStorageLocationByObjectInstanceId for object instance ID: %v", copyFromTo.ObjectInstance.Id)
	}

	connection := models.Connection{}
	err = json.Unmarshal([]byte(copyFromTo.LocationCopyTo.Connection), &connection)
	if err != nil {
		c.Logger.Error().Msgf("error mapping json")
		return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	daLogger := zLogger.NewZWrapper(c.Logger)
	vfsConfig, err := config.LoadVfsConfig(copyFromTo.LocationCopyTo.Connection, storageLocationToCopyFrom.Connection)
	if err != nil {
		c.Logger.Error().Msgf("error mapping json for storage location connection field: %v", err)
		return nil, errors.Wrapf(err, "error mapping json for storage location connection field")
	}
	vfs, err := vfsrw.NewFS(vfsConfig, daLogger)
	if err != nil {
		c.Logger.Warn().Msgf("cannot create vfs: %v", err)
		return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	path := connection.Folder + storagePartition.Alias + "/" + filepath.Base(copyFromTo.ObjectInstance.Path)
	sourceFP, err := vfs.Open(copyFromTo.ObjectInstance.Path)

	if err == nil {
		storagePartition.CurrentSize += copyFromTo.ObjectInstance.Size
		storagePartition.CurrentObjects++
		objectInstance := &pb.ObjectInstance{Path: path, Status: "new", ObjectId: copyFromTo.ObjectInstance.ObjectId, StoragePartitionId: storagePartition.Id, Size: copyFromTo.ObjectInstance.Size}
		_, err = c.ClientStorageHandlerHandler.CreateObjectInstance(ctx, objectInstance)
		if err != nil {
			c.Logger.Error().Msgf("Could not create objectInstance for object with ID: %v", copyFromTo.ObjectInstance.ObjectId)
			return &pb.NoParam{}, errors.Wrapf(err, "Could not create objectInstance for object with ID: %v", copyFromTo.ObjectInstance.ObjectId)
		}
		_, err = c.ClientStorageHandlerHandler.UpdateStoragePartition(ctx, storagePartition)
		if err != nil {
			c.Logger.Error().Msgf("Could not update storage partition with ID: %v", storagePartition.Id)
			return &pb.NoParam{}, errors.Wrapf(err, "Could not update storage partition with ID: %v", storagePartition.Id)
		}

		targetFP, err := writefs.Create(vfs, path)
		if err != nil {
			c.Logger.Error().Msgf("cannot create target for path '%v': %s", path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "cannot create target for path '%v': %s", path, err)
		}
		defer func() {
			if err := targetFP.Close(); err != nil {
				c.Logger.Error().Msgf("cannot close target: %v", err)
			}
		}()
		csWriter, err := checksum.NewChecksumWriter(
			[]checksum.DigestAlgorithm{checksum.DigestSHA512},
			targetFP,
		)
		if err != nil {
			c.Logger.Error().Msgf("cannot create checksum writer for file '%s%s': %v", vfs, path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "cannot create checksum writer for file '%v%v': %s", vfs, path, err)
		}
		_, err = io.Copy(csWriter, sourceFP)
		if err != nil {
			if err := csWriter.Close(); err != nil {
				c.Logger.Error().Msgf("cannot close checksum writer: %v", err)
				return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			c.Logger.Error().Msgf("error writing file to path '%v%v': %s", vfs, path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "error writing file to path '%v%v': %s", vfs, path, err)
		}
		if err := csWriter.Close(); err != nil {
			c.Logger.Error().Msgf("cannot close checksum writer: %v", err)
			return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
		}
		sourceFP.Close()
		vfs.Close()

	} else {
		c.Logger.Error().Msgf("error opening file with path %v: %s", copyFromTo.ObjectInstance.Path, err)
		return nil, errors.Wrapf(err, "error opening file with path %v", copyFromTo.ObjectInstance.Path)
	}
	return &pb.NoParam{}, nil
}
