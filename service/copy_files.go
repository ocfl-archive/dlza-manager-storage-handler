package service

import (
	"context"
	"io"
	"io/fs"

	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/je4/utils/v2/pkg/zLogger"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func StoringFiles(clientStorageHandlerHandler pbHandler.StorageHandlerHandlerServiceClient, ctx context.Context, objectWithCollectionAliasAndPathAndFiles *pb.IncomingOrder, partitionId string, severalObjects string, logger zLogger.ZLogger) (*pb.Status, error) {

	stream, err := clientStorageHandlerHandler.SaveAllTableObjectsAfterCopyingStream(ctx)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot SaveAllTableObjectsAfterCopying for collection with alias: %v and path: %v", objectWithCollectionAliasAndPathAndFiles.CollectionAlias, objectWithCollectionAliasAndPathAndFiles.FilePath)
	}

	instanceWithPartitionAndObjectWithFile := &pb.InstanceWithPartitionAndObjectWithFile{}
	if objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Binary && severalObjects != "1" { // only if json does not contain files. "1" means that it contains
		instanceWithPartitionAndObjectWithFile.Object = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object
		instanceWithPartitionAndObjectWithFile.StoragePartition = &pb.StoragePartition{Id: partitionId}
		instanceWithPartitionAndObjectWithFile.ObjectInstance = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.ObjectInstance
		instanceWithPartitionAndObjectWithFile.CollectionAlias = objectWithCollectionAliasAndPathAndFiles.CollectionAlias
		instanceWithPartitionAndObjectWithFile.NewVersion = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.NewVersion
		if err := stream.Send(instanceWithPartitionAndObjectWithFile); err != nil {
			return &pb.Status{Ok: false}, errors.Wrapf(err, "Could store all table objects for collection: %v", objectWithCollectionAliasAndPathAndFiles.CollectionAlias)
		}
	} else {
		for i, objectAndFile := range objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Files {
			if i == 0 {
				instanceWithPartitionAndObjectWithFile.Object = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object
				instanceWithPartitionAndObjectWithFile.StoragePartition = &pb.StoragePartition{Id: partitionId}
				instanceWithPartitionAndObjectWithFile.ObjectInstance = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.ObjectInstance
				instanceWithPartitionAndObjectWithFile.CollectionAlias = objectWithCollectionAliasAndPathAndFiles.CollectionAlias
				instanceWithPartitionAndObjectWithFile.NewVersion = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.NewVersion
			}
			instanceWithPartitionAndObjectWithFile.File = objectAndFile
			if err := stream.Send(instanceWithPartitionAndObjectWithFile); err != nil {
				return &pb.Status{Ok: false}, errors.Wrapf(err, "Could store all table objects for collection: %v", objectWithCollectionAliasAndPathAndFiles.CollectionAlias)
			}
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "Could store all table objects: %v, CloseAndRecv failed", objectWithCollectionAliasAndPathAndFiles.CollectionAlias)
	}

	return &pb.Status{Ok: true}, nil
}

func DeleteTemporaryFiles(filePaths []string, vfs fs.FS, logger zLogger.ZLogger) (*pb.Status, error) {
	for _, filePath := range filePaths {
		if err := writefs.Remove(vfs, filePath); err != nil {
			logger.Error().Msgf("error deleting file to '%s': %v", filePath, err)
			return &pb.Status{Ok: false}, grpcstatus.Errorf(codes.Internal, "error deleting file '%s': %v", filePath, err)
		}
	}
	return &pb.Status{Ok: true}, nil
}
