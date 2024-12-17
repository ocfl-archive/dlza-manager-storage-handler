package server

import (
	"context"
	"emperror.dev/errors"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"io"
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
