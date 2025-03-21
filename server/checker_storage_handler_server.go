package server

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/dlza-manager/models"
	"io"
	"io/fs"
	"path/filepath"
)

type CheckerStorageHandlerServer struct {
	storageHandlerPb.UnimplementedCheckerStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
	Logger                      zLogger.ZLogger
	Vfs                         fs.FS
}

func (c *CheckerStorageHandlerServer) GetObjectInstanceChecksum(ctx context.Context, objectInstance *pb.ObjectInstance) (*pb.Id, error) {

	sourceFP, err := c.Vfs.Open(objectInstance.Path)
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
		sourceFP.Close()
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
	sourceFP.Close()
	return &pb.Id{Id: checksums[checksum.DigestSHA512]}, nil
}

func (c *CheckerStorageHandlerServer) CopyArchiveTo(ctx context.Context, copyFromTo *pb.CopyFromTo) (*pb.NoParam, error) {
	storagePartition, err := c.ClientStorageHandlerHandler.GetStoragePartitionForLocation(ctx, &pb.SizeAndId{Size: copyFromTo.ObjectInstance.Size, Id: copyFromTo.LocationCopyTo.Id})
	if err != nil {
		c.Logger.Error().Msgf("cannot get storagePartition for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
		return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	connection := models.Connection{}
	err = json.Unmarshal([]byte(copyFromTo.LocationCopyTo.Connection), &connection)
	if err != nil {
		c.Logger.Error().Msgf("error mapping json")
		return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	path := connection.Folder + storagePartition.Alias + "/" + filepath.Base(copyFromTo.ObjectInstance.Path)
	sourceFP, err := c.Vfs.Open(copyFromTo.ObjectInstance.Path)

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

		targetFP, err := writefs.Create(c.Vfs, path)
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
			c.Logger.Error().Msgf("cannot create checksum writer for file '%s%s': %v", c.Vfs, path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "cannot create checksum writer for file '%v%v': %s", c.Vfs, path, err)
		}
		_, err = io.Copy(csWriter, sourceFP)
		if err != nil {
			if err := csWriter.Close(); err != nil {
				c.Logger.Error().Msgf("cannot close checksum writer: %v", err)
				return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			c.Logger.Error().Msgf("error writing file to path '%v%v': %s", c.Vfs, path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "error writing file to path '%v%v': %s", c.Vfs, path, err)
		}
		if err := csWriter.Close(); err != nil {
			c.Logger.Error().Msgf("cannot close checksum writer: %v", err)
			return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
		}
		sourceFP.Close()
	} else {
		c.Logger.Error().Msgf("error opening file with path %v: %s", copyFromTo.ObjectInstance.Path, err)
		return nil, errors.Wrapf(err, "error opening file with path %v", copyFromTo.ObjectInstance.Path)
	}
	return &pb.NoParam{}, nil
}
