package server

import (
	"context"
	zw "github.com/je4/utils/v2/pkg/zLogger"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/pkg/errors"
)

type UploaderStorageHandlerServer struct {
	storageHandlerPb.UnimplementedUploaderStorageHandlerServiceServer
	ClientStorageHandlerHandler pbHandler.StorageHandlerHandlerServiceClient
	ConfigObj                   config.Config
	Logger                      zw.ZWrapper
}

func (u *UploaderStorageHandlerServer) CopyFile(ctx context.Context, incomingOrder *pb.IncomingOrder) (*pb.Status, error) {
	_, err := u.ClientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: incomingOrder.StatusId, Status: "archiving"})
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot set status to copy file for collection '%s'", incomingOrder.CollectionAlias)
	}
	status, err := service.CopyFiles(u.ClientStorageHandlerHandler, ctx, incomingOrder, u.ConfigObj, u.Logger)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot copy file for collection '%s'", incomingOrder.CollectionAlias)
	}
	_, err = u.ClientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: incomingOrder.StatusId, Status: "archived"})
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot set status to copy file for collection '%s'", incomingOrder.CollectionAlias)
	}
	_, err = service.DeleteTemporaryFiles(incomingOrder, u.ConfigObj, u.Logger)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot delete temporary files for collection '%s'", incomingOrder.CollectionAlias)
	}
	return status, nil
}
