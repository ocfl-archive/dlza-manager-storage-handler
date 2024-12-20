package service

import (
	"context"
	"encoding/json"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
	"github.com/je4/filesystem/v2/pkg/writefs"
	"github.com/je4/utils/v2/pkg/checksum"
	config2 "github.com/je4/utils/v2/pkg/config"
	zw "github.com/je4/utils/v2/pkg/zLogger"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/pkg/errors"
	"io"
	"maps"
)

func CopyFiles(clientStorageHandlerHandler pbHandler.StorageHandlerHandlerServiceClient, ctx context.Context, objectWithCollectionAliasAndPathAndFiles []*pb.ObjectAndFile, cfg config.Config, daLogger zw.ZWrapper) (*pb.Status, error) {

	storageLocations, err := clientStorageHandlerHandler.GetStorageLocationsByCollectionAlias(ctx, &pb.CollectionAlias{CollectionAlias: objectWithCollectionAliasAndPathAndFiles[0].CollectionAlias})

	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot get storageLocations for collection: %v", objectWithCollectionAliasAndPathAndFiles[0].CollectionAlias)
	}
	configObj, err := models.LoadStorageHandlerConfig(storageLocations.StorageLocations)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot load StorageHandler config: %v", err)
	}

	tempVfsMap := getVfsTempMap(cfg)

	maps.Copy(tempVfsMap, configObj.VFS)
	configObj.VFS = tempVfsMap

	vfs, err := vfsrw.NewFS(configObj.VFS, daLogger)
	if err != nil {
		daLogger.Errorf("cannot create vfs: %v", err)
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot create vfs: %v", err)
	}

	var storageLocation *pb.StorageLocation
	for _, storageLocationItem := range storageLocations.StorageLocations {
		if storageLocationItem.FillFirst {
			storageLocation = storageLocationItem
		}
	}

	storagePartition, err := clientStorageHandlerHandler.GetStoragePartitionForLocation(ctx, &pb.SizeAndId{Size: objectWithCollectionAliasAndPathAndFiles[0].Object.Size, Id: storageLocation.Id})
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", storageLocation.Alias)
	}

	connection := models.Connection{}
	err = json.Unmarshal([]byte(storageLocation.Connection), &connection)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "error mapping storageLocation json for storageLocation ID: %v", storageLocation.Id)
	}

	path := connection.Folder + storagePartition.Alias + "/" + objectWithCollectionAliasAndPathAndFiles[0].FileName

	objectInstance := &pb.ObjectInstance{Path: path, Status: "new", StoragePartitionId: storagePartition.Id, Size: objectWithCollectionAliasAndPathAndFiles[0].Object.Size}
	storagePartition.CurrentSize += objectWithCollectionAliasAndPathAndFiles[0].Object.Size
	storagePartition.CurrentObjects++

	stream, err := clientStorageHandlerHandler.SaveAllTableObjectsAfterCopyingStream(ctx)

	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot SaveAllTableObjectsAfterCopying for collection with alias: %v and path: %v", objectWithCollectionAliasAndPathAndFiles[0].CollectionAlias, path)
	}
	for i, objectAndFile := range objectWithCollectionAliasAndPathAndFiles {
		instanceWithPartitionAndObjectWithFile := &pb.InstanceWithPartitionAndObjectWithFile{}
		if i == 0 {
			instanceWithPartitionAndObjectWithFile.Object = objectWithCollectionAliasAndPathAndFiles[i].Object
			instanceWithPartitionAndObjectWithFile.StoragePartition = storagePartition
			instanceWithPartitionAndObjectWithFile.ObjectInstance = objectInstance
			instanceWithPartitionAndObjectWithFile.CollectionAlias = objectWithCollectionAliasAndPathAndFiles[i].CollectionAlias
		}
		instanceWithPartitionAndObjectWithFile.File = objectAndFile.File
		if err := stream.Send(instanceWithPartitionAndObjectWithFile); err != nil {
			return &pb.Status{Ok: false}, errors.Wrapf(err, "Could store all table objects for collection: %v", objectWithCollectionAliasAndPathAndFiles[0].CollectionAlias)
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "Could store all table objects: %v, CloseAndRecv failed", objectWithCollectionAliasAndPathAndFiles[0].CollectionAlias)
	}
	err = func() error {

		sourceFP, err := vfs.Open(objectWithCollectionAliasAndPathAndFiles[0].FilePath)
		if err != nil {
			daLogger.Errorf("cannot read file '%s': %v", objectWithCollectionAliasAndPathAndFiles[0].FilePath, err)
			return errors.Wrapf(err, "cannot read file '%s': %v", objectWithCollectionAliasAndPathAndFiles[0].FilePath, err)
		}
		defer func() {
			if err := sourceFP.Close(); err != nil {
				daLogger.Errorf("cannot close source: %v", err)
			}
		}()

		targetFP, err := writefs.Create(vfs, path)
		if err != nil {
			return errors.Wrapf(err, "cannot create file '%s%s': %v", vfs, path, err)
		}
		defer func() {
			if err := targetFP.Close(); err != nil {
				daLogger.Errorf("cannot close target: %v", err)
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
			daLogger.Errorf("error writing file")
			if err := csWriter.Close(); err != nil {
				daLogger.Errorf("cannot close checksum writer: %v", err)
				return errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			return errors.Wrapf(err, "error writing file: %v", objectWithCollectionAliasAndPathAndFiles[0].FilePath)
		}
		if _size != objectWithCollectionAliasAndPathAndFiles[0].Object.Size {
			if err := csWriter.Close(); err != nil {
				daLogger.Errorf("cannot close checksum writer: %v", err)
				return errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			return errors.Wrapf(err, "size should be the same: '%v' != %v", _size, objectWithCollectionAliasAndPathAndFiles[0].Object.Size)
		}

		if err := csWriter.Close(); err != nil {
			daLogger.Errorf("cannot close checksum writer: %v", err)
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

	_, err = clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: objectWithCollectionAliasAndPathAndFiles[0].StatusId, Status: "zip was copied"})
	if err != nil {
		daLogger.Warningf("could not AlterStatus with status id %s:  to zip was copied", objectWithCollectionAliasAndPathAndFiles[0].StatusId)
	}

	return &pb.Status{Ok: true}, nil
}

func DeleteTemporaryFiles(objectAndFile *pb.ObjectAndFile, cfg config.Config, daLogger zw.ZWrapper) (*pb.Status, error) {
	tempVfsMap := getVfsTempMap(cfg)
	vfs, err := vfsrw.NewFS(tempVfsMap, daLogger)
	if err != nil {
		daLogger.Errorf("cannot create vfs: %v", err)
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot create vfs: %v", err)
	}

	if err := writefs.Remove(vfs, objectAndFile.FilePath); err != nil {
		daLogger.Errorf("error deleting file to '%s': %v", objectAndFile.FilePath, err)
		return &pb.Status{Ok: false}, errors.Wrapf(err, "error writing file to '%s': %v", objectAndFile.FilePath, err)
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
