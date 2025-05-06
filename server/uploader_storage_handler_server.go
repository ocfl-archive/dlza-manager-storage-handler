package server

import (
	"context"
	"emperror.dev/errors"
	zw "github.com/je4/utils/v2/pkg/zLogger"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"io"
)

type UploaderStorageHandlerServer struct {
	storageHandlerPb.UnimplementedUploaderStorageHandlerServiceServer
	ClientStorageHandlerHandler pbHandler.StorageHandlerHandlerServiceClient
	ConfigObj                   config.Config
	Logger                      zw.ZWrapper
}

func (u *UploaderStorageHandlerServer) CopyFileStream(stream storageHandlerPb.UploaderStorageHandlerService_CopyFileStreamServer) error {
	var objectAndFiles []*pb.ObjectAndFile
	for {
		file, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		objectAndFiles = append(objectAndFiles, file)
	}
	_, err := u.ClientStorageHandlerHandler.AlterStatus(context.Background(), &pb.StatusObject{Id: objectAndFiles[0].StatusId, Status: "archiving"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", objectAndFiles[0].CollectionAlias)
	}
	_, err = service.CopyFiles(u.ClientStorageHandlerHandler, context.Background(), objectAndFiles, u.ConfigObj, u.Logger)
	if err != nil {
		return errors.Wrapf(err, "cannot copy file for collection '%s'", objectAndFiles[0].CollectionAlias)
	}
	_, err = u.ClientStorageHandlerHandler.AlterStatus(context.Background(), &pb.StatusObject{Id: objectAndFiles[0].StatusId, Status: "archived"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", objectAndFiles[0].CollectionAlias)
	}
	_, err = service.DeleteTemporaryFiles(objectAndFiles[0], u.ConfigObj, u.Logger)
	if err != nil {
		return errors.Wrapf(err, "cannot delete temporary files for collection '%s'", objectAndFiles[0].CollectionAlias)
	}
	return nil
}
