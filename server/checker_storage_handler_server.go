package server

import (
	"context"
	"emperror.dev/errors"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"io"
	"io/fs"
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
