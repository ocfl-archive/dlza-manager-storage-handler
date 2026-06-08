package server

import (
	"context"
	"io/fs"

	"emperror.dev/errors"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/filesystem/pkg/writefs"
)

type ClerkStorageHandlerServer struct {
	storageHandlerPb.UnimplementedClerkStorageHandlerServiceServer
	Vfs fs.FS
}

func (c *ClerkStorageHandlerServer) CreateFolder(ctx context.Context, path *pb.Id) (*pb.NoParam, error) {
	if err := writefs.MkDir(c.Vfs, path.Id); err != nil {
		return &pb.NoParam{}, errors.Wrapf(err, "error creating folder with path: %s", path.Id)
	}
	return &pb.NoParam{}, nil
}
