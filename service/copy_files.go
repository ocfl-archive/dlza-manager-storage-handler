package service

import (
	"context"
	"encoding/json"
	"github.com/je4/filesystem/v3/pkg/vfsrw"
	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/je4/utils/v2/pkg/checksum"
	config2 "github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/zLogger"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/pkg/errors"
	"io"
	"io/fs"
)

func CopyFiles(clientStorageHandlerHandler pbHandler.StorageHandlerHandlerServiceClient, ctx context.Context, objectWithCollectionAliasAndPathAndFiles *pb.IncomingOrder, vfs fs.FS, logger zLogger.ZLogger) (*pb.Status, error) {

	storageLocations, err := clientStorageHandlerHandler.GetStorageLocationsByCollectionAlias(ctx, &pb.CollectionAlias{CollectionAlias: objectWithCollectionAliasAndPathAndFiles.CollectionAlias})

	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot get storageLocations for collection: %v", objectWithCollectionAliasAndPathAndFiles.CollectionAlias)
	}

	var storageLocation *pb.StorageLocation
	for _, storageLocationItem := range storageLocations.StorageLocations {
		if storageLocationItem.FillFirst {
			storageLocation = storageLocationItem
		}
	}

	storagePartition, err := clientStorageHandlerHandler.GetStoragePartitionForLocation(ctx, &pb.SizeAndId{Size: objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Size, Id: storageLocation.Id, Object: objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object})
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", storageLocation.Alias)
	}

	connection := models.Connection{}
	err = json.Unmarshal([]byte(storageLocation.Connection), &connection)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "error mapping storageLocation json for storageLocation ID: %v", storageLocation.Id)
	}

	path := connection.Folder + storagePartition.Alias + "/" + objectWithCollectionAliasAndPathAndFiles.FileName

	objectInstance := &pb.ObjectInstance{Path: path, Status: "new", StoragePartitionId: storagePartition.Id, Size: objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Size}
	storagePartition.CurrentSize += objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Size
	storagePartition.CurrentObjects++

	stream, err := clientStorageHandlerHandler.SaveAllTableObjectsAfterCopyingStream(ctx)

	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot SaveAllTableObjectsAfterCopying for collection with alias: %v and path: %v", objectWithCollectionAliasAndPathAndFiles.CollectionAlias, path)
	}
	instanceWithPartitionAndObjectWithFile := &pb.InstanceWithPartitionAndObjectWithFile{}
	if objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Binary {
		instanceWithPartitionAndObjectWithFile.Object = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object
		instanceWithPartitionAndObjectWithFile.StoragePartition = storagePartition
		instanceWithPartitionAndObjectWithFile.ObjectInstance = objectInstance
		instanceWithPartitionAndObjectWithFile.CollectionAlias = objectWithCollectionAliasAndPathAndFiles.CollectionAlias
		if err := stream.Send(instanceWithPartitionAndObjectWithFile); err != nil {
			return &pb.Status{Ok: false}, errors.Wrapf(err, "Could store all table objects for collection: %v", objectWithCollectionAliasAndPathAndFiles.CollectionAlias)
		}
	} else {
		for i, objectAndFile := range objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Files {
			if i == 0 {
				instanceWithPartitionAndObjectWithFile.Object = objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object
				instanceWithPartitionAndObjectWithFile.StoragePartition = storagePartition
				instanceWithPartitionAndObjectWithFile.ObjectInstance = objectInstance
				instanceWithPartitionAndObjectWithFile.CollectionAlias = objectWithCollectionAliasAndPathAndFiles.CollectionAlias
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
	err = func() error {

		sourceFP, err := vfs.Open(objectWithCollectionAliasAndPathAndFiles.FilePath)
		if err != nil {
			logger.Error().Msgf("cannot read file '%s': %v", objectWithCollectionAliasAndPathAndFiles.FilePath, err)
			return errors.Wrapf(err, "cannot read file '%s': %v", objectWithCollectionAliasAndPathAndFiles.FilePath, err)
		}
		defer func() {
			if err := sourceFP.Close(); err != nil {
				logger.Error().Msgf("cannot close source: %v", err)
			}
		}()
		targetFP, err := writefs.Create(vfs, path)
		if err != nil {
			return errors.Wrapf(err, "cannot create file '%s%s': %v", vfs, path, err)
		}
		defer func() {
			if err := targetFP.Close(); err != nil {
				logger.Error().Msgf("cannot close target: %v", err)
			}
		}()
		csWriter, err := checksum.NewChecksumWriter(
			[]checksum.DigestAlgorithm{checksum.DigestSHA512},
			targetFP,
		)
		if err != nil {
			return errors.Wrapf(err, "cannot create checksum writer for file '%s%s': %v", vfs, path, err)
		}

		_size, err := io.Copy(csWriter, sourceFP)
		if err != nil {
			logger.Error().Msgf("error writing file")
			if err := csWriter.Close(); err != nil {
				logger.Error().Msgf("cannot close checksum writer: %v", err)
				return errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			return errors.Wrapf(err, "error writing file: %v", objectWithCollectionAliasAndPathAndFiles.FilePath)
		}
		if _size != objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Size {
			if err := csWriter.Close(); err != nil {
				logger.Error().Msgf("cannot close checksum writer: %v", err)
				return errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			return errors.Wrapf(err, "size should be the same: '%v' != %v", _size, objectWithCollectionAliasAndPathAndFiles.ObjectAndFiles.Object.Size)
		}

		if err := csWriter.Close(); err != nil {
			logger.Error().Msgf("cannot close checksum writer: %v", err)
			return errors.Wrapf(err, "cannot close checksum writer: %v", err)
		}
		/*
			checksums, err := csWriter.GetChecksums()
			if err != nil {
				daLogger.Errorf("cannot get checksum: %v", err)
				return errors.Wrapf(err, "cannot get checksum: %v", err)
			}
		*/

		return nil
	}()
	if err != nil {
		return &pb.Status{Ok: false}, err
	}

	_, err = clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: objectWithCollectionAliasAndPathAndFiles.StatusId, Status: "zip was copied"})
	if err != nil {
		logger.Warn().Msgf("could not AlterStatus with status id %s:  to zip was copied", objectWithCollectionAliasAndPathAndFiles.StatusId)
	}

	return &pb.Status{Ok: true}, nil
}

func DeleteTemporaryFiles(filePath string, cfg config.Config, logger zLogger.ZLogger) (*pb.Status, error) {
	tempVfsMap := getVfsTempMap(cfg)
	vfs, err := vfsrw.NewFS(tempVfsMap, logger)
	if err != nil {
		logger.Error().Msgf("cannot create vfs: %v", err)
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot create vfs: %v", err)
	}

	if err := writefs.Remove(vfs, filePath); err != nil {
		logger.Error().Msgf("error deleting file '%s': %v", filePath, err)
		return &pb.Status{Ok: false}, errors.Wrapf(err, "error deleting file to '%s': %v", filePath, err)
	}

	return &pb.Status{Ok: true}, nil
}

func getVfsTempMap(cfg config.Config) map[string]*vfsrw.VFS {
	vfsTemp := vfsrw.VFS{
		Type: cfg.S3TempStorage.Type,
		Name: cfg.S3TempStorage.Name,
		S3: &vfsrw.S3{
			AccessKeyID:     config2.EnvString(cfg.S3TempStorage.Key),
			SecretAccessKey: config2.EnvString(cfg.S3TempStorage.Secret),
			Endpoint:        config2.EnvString(cfg.S3TempStorage.Url),
			Region:          "us-east-1",
			UseSSL:          true,
			Debug:           cfg.S3TempStorage.Debug,
			CAPEM:           cfg.S3TempStorage.CAPEM,
		},
	}

	tempVfsMap := make(map[string]*vfsrw.VFS)
	tempVfsMap[cfg.S3TempStorage.Name] = &vfsTemp
	return tempVfsMap
}
