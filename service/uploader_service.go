package service

import (
	"context"
	"emperror.dev/errors"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"time"
)

type UploaderService struct {
	StorageHandlerHandlerServiceClient handlerPb.StorageHandlerHandlerServiceClient
	Logger                             *zLogger.ZLogger
}

func (u *UploaderService) TenantHasAccess(key string, collection string) (bool, error) {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	status, err := u.StorageHandlerHandlerServiceClient.TenantHasAccess(ctx, &pb.UploaderAccessObject{Key: key, Collection: collection})
	if err != nil {
		return false, errors.Wrapf(err, "could not get tenant access status for tenant with key: %v", key)
	}
	return status.Ok, nil
}
