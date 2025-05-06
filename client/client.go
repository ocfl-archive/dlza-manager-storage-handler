package client

import (
	"emperror.dev/errors"
	pb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	"google.golang.org/grpc"
	"io"
)

const maxMsgSize = 1024 * 1024 * 12

func NewStorageHandlerClerkClient(target string, opt grpc.DialOption) (pb.ClerkStorageHandlerServiceClient, io.Closer, error) {
	connection, err := grpc.NewClient(target, opt, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return pb.NewClerkStorageHandlerServiceClient(connection), connection, nil
}

func NewUploaderStorageHandlerClient(target string, opt grpc.DialOption) (pb.UploaderStorageHandlerServiceClient, io.Closer, error) {
	connection, err := grpc.NewClient(target, opt, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return pb.NewUploaderStorageHandlerServiceClient(connection), connection, nil
}

func NewDispatcherStorageHandlerClient(target string, opt grpc.DialOption) (pb.DispatcherStorageHandlerServiceClient, io.Closer, error) {
	connection, err := grpc.NewClient(target, opt, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return pb.NewDispatcherStorageHandlerServiceClient(connection), connection, nil
}
