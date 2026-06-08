package server

import (
	"context"
	"io"
	"io/fs"

	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DispatcherStorageHandlerServer struct {
	storageHandlerPb.UnimplementedDispatcherStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
	Logger                      zLogger.ZLogger
	Vfs                         fs.FS
}

func (d *DispatcherStorageHandlerServer) CopyArchiveTo(ctx context.Context, copyFromTo *pb.CopyFromTo) (*pb.NoParam, error) {

	sourceFP, err := d.Vfs.Open(copyFromTo.CopyFrom)
	if err != nil {
		d.Logger.Error().Msgf("could not open file with path %s, err: %s", copyFromTo.CopyFrom, err)
		return &pb.NoParam{}, errors.Wrapf(err, "could not open file with path %s", copyFromTo.CopyFrom)
	}

	targetFP, err := writefs.Create(d.Vfs, copyFromTo.CopyTo)
	if err != nil {
		d.Logger.Error().Msgf("cannot create target for path '%v': %s", copyFromTo.CopyTo, err)
		if err := sourceFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
		}
		return &pb.NoParam{}, errors.Wrapf(err, "cannot create target for path '%v': %s", copyFromTo.CopyTo, err)
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
		d.Logger.Error().Msgf("cannot create checksum writer for file '%s%s': %v", d.Vfs, copyFromTo.CopyTo, err)
		if err := sourceFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
		}
		return &pb.NoParam{}, errors.Wrapf(err, "cannot create checksum writer for file '%v%v': %s", d.Vfs, copyFromTo.CopyTo, err)
	}
	_, err = io.Copy(csWriter, sourceFP)
	if err != nil {
		if err := csWriter.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
			return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
		}
		d.Logger.Error().Msgf("error writing file to path '%v%v': %s", d.Vfs, copyFromTo.CopyTo, err)
		if err := sourceFP.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close sourceFP: %v", err)
		}
		return &pb.NoParam{}, errors.Wrapf(err, "error writing file to path '%v%v': %s", d.Vfs, copyFromTo.CopyTo, err)
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

func (d *DispatcherStorageHandlerServer) ConnectionCheck(ctx context.Context, empty *emptypb.Empty) (*pb.Id, error) {
	return &pb.Id{Id: "Storage Handler -> Dispatcher connection checked"}, nil
}
