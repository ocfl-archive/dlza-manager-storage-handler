package server

import (
	"context"
	"encoding/json"
	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"path/filepath"
)

type DispatcherStorageHandlerServer struct {
	storageHandlerPb.UnimplementedDispatcherStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
	Logger                      zLogger.ZLogger
	Vfs                         fs.FS
}

func (d *DispatcherStorageHandlerServer) CopyArchiveTo(ctx context.Context, copyFromTo *pb.CopyFromTo) (*pb.NoParam, error) {

	connection := models.Connection{}
	err := json.Unmarshal([]byte(copyFromTo.LocationCopyTo.Connection), &connection)
	if err != nil {
		d.Logger.Error().Msgf("error mapping json")
		return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	path := connection.Folder + copyFromTo.StoragePartition.Alias + "/" + filepath.Base(copyFromTo.ObjectInstance.Path)
	sourceFP, err := d.Vfs.Open(copyFromTo.ObjectInstance.Path)
	if err != nil {
		d.Logger.Error().Msgf("could not open file with path %s, err: %s", copyFromTo.ObjectInstance.Path, err)
		return &pb.NoParam{}, errors.Wrapf(err, "could not open file with path %s", copyFromTo.ObjectInstance.Path)
	}
	objectInstance := &pb.ObjectInstance{Path: path, Status: "new", ObjectId: copyFromTo.ObjectInstance.ObjectId, StoragePartitionId: copyFromTo.StoragePartition.Id, Size: copyFromTo.ObjectInstance.Size}
	_, err = d.ClientStorageHandlerHandler.CreateObjectInstance(ctx, objectInstance)
	if err != nil {
		d.Logger.Error().Msgf("Could not create objectInstance for object with ID: %v", copyFromTo.ObjectInstance.ObjectId)
		return &pb.NoParam{}, errors.Wrapf(err, "Could not create objectInstance for object with ID: %v", copyFromTo.ObjectInstance.ObjectId)
	}

	targetFP, err := writefs.Create(d.Vfs, path)
	if err != nil {
		d.Logger.Error().Msgf("cannot create target for path '%v': %s", path, err)
		if err := sourceFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
		}
		return &pb.NoParam{}, errors.Wrapf(err, "cannot create target for path '%v': %s", path, err)
	}
	defer func() {
		if err := targetFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close target: %v", err)
		}
	}()
	csWriter, err := checksum.NewChecksumWriter(
		[]checksum.DigestAlgorithm{checksum.DigestSHA512},
		targetFP,
	)
	if err != nil {
		d.Logger.Error().Msgf("cannot create checksum writer for file '%s%s': %v", d.Vfs, path, err)
		if err := sourceFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
		}
		return &pb.NoParam{}, errors.Wrapf(err, "cannot create checksum writer for file '%v%v': %s", d.Vfs, path, err)
	}
	_, err = io.Copy(csWriter, sourceFP)
	if err != nil {
		if err := csWriter.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
			return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
		}
		d.Logger.Error().Msgf("error writing file to path '%v%v': %s", d.Vfs, path, err)
		if err := sourceFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
		}
		return &pb.NoParam{}, errors.Wrapf(err, "error writing file to path '%v%v': %s", d.Vfs, path, err)
	}
	if err := csWriter.Close(); err != nil {
		d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
		return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
	}
	if err := sourceFP.Close(); err != nil {
		d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
	}

	return &pb.NoParam{}, nil
}
