package server

import (
	"context"
	"encoding/json"
	"github.com/je4/utils/v2/pkg/checksum"
	zw "github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"io"
	"path/filepath"

	"github.com/je4/filesystem/v2/pkg/vfsrw"
	"github.com/je4/filesystem/v2/pkg/writefs"
	"github.com/pkg/errors"
)

type DispatcherStorageHandlerServer struct {
	storageHandlerPb.UnimplementedDispatcherStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
	Logger                      zw.ZWrapper
}

func (d *DispatcherStorageHandlerServer) ChangeQualityForCollections(ctx context.Context, collectionAliases *pb.CollectionAliases) (*pb.NoParam, error) {

	for _, collectionAlias := range collectionAliases.CollectionAliases {
		//get cheapest storage locations
		storageLocationsPb, err := d.ClientStorageHandlerHandler.GetStorageLocationsByCollectionAlias(ctx, collectionAlias)
		if err != nil {
			return &pb.NoParam{}, errors.Wrapf(err, "cannot get storageLocations for collection: %v", collectionAlias)
		}
		if storageLocationsPb.StorageLocations == nil {
			d.Logger.Warningf("The collection " + collectionAlias.CollectionAlias + " doesn't have enough storage locations")
			continue
		}

		ObjectsPb, err := d.ClientStorageHandlerHandler.GetObjectsByCollectionAlias(ctx, collectionAlias)
		if err != nil {
			return &pb.NoParam{}, errors.Wrapf(err, "cannot get objects for collection: %v", collectionAlias)
		}

		for _, objectPb := range ObjectsPb.Objects {

			currentStorageLocationsPb, err := d.ClientStorageHandlerHandler.GetStorageLocationsByObjectId(ctx, &pb.Id{Id: objectPb.Id})
			if err != nil {
				return &pb.NoParam{}, errors.Wrapf(err, "cannot GetCurrentStorageLocationsByCollectionAlias: %v", err)
			}

			storageLocationsToCopyIn := service.GetStorageLocationsToCopyIn(storageLocationsPb, currentStorageLocationsPb)
			storageLocationsToDeleteFrom := service.GetStorageLocationsToDeleteFrom(storageLocationsPb, currentStorageLocationsPb)
			if len(storageLocationsToCopyIn.StorageLocations) == 0 && len(storageLocationsToCopyIn.StorageLocations) == 0 {
				continue
			}
			config, err := models.LoadStorageHandlerConfig(append(currentStorageLocationsPb.StorageLocations, storageLocationsToCopyIn.StorageLocations...))
			if err != nil {
				return &pb.NoParam{}, errors.Wrapf(err, "cannot load storage-handler config: %v", err)
			}

			vfs, err := vfsrw.NewFS(config.VFS, d.Logger)
			if err != nil {
				d.Logger.Warningf("cannot create vfs: %v", err)
				continue
			}

			ObjectInstancesPb, err := d.ClientStorageHandlerHandler.GetObjectsInstancesByObjectId(ctx, &pb.Id{Id: objectPb.Id})
			if err != nil {
				return &pb.NoParam{}, errors.Wrapf(err, "cannot get object instances for collection: %v", collectionAlias)
			}
			//ToDo Find a better way to chose the path to copy from
			pathToCopyFrom := ObjectInstancesPb.ObjectInstances[0].Path

			for _, storageLocation := range storageLocationsToCopyIn.StorageLocations {
				storagePartition, err := d.ClientStorageHandlerHandler.GetStoragePartitionForLocation(ctx, &pb.SizeAndId{Size: int64(objectPb.Size), Id: storageLocation.Id})
				if err != nil {
					d.Logger.Errorf("cannot get storagePartition for storageLocation: %v", storageLocation.Alias)
					return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", storageLocation.Alias)
				}
				connection := models.Connection{}
				err = json.Unmarshal([]byte(storageLocation.Connection), &connection)
				if err != nil {
					d.Logger.Errorf("error mapping json")
					return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", storageLocation.Alias)
				}

				path := connection.Folder + storagePartition.Alias + "/" + filepath.Base(pathToCopyFrom)
				err = func() error {
					sourceFP, err := vfs.Open(pathToCopyFrom)
					if err != nil {
						d.Logger.Errorf("cannot open file '%s': %v", pathToCopyFrom, err)
						return errors.Wrapf(err, "cannot open file '%s': %v", pathToCopyFrom, err)
					}
					defer func() {
						if err := sourceFP.Close(); err != nil {
							d.Logger.Errorf("cannot close source: %v", err)
						}
					}()

					targetFP, err := writefs.Create(vfs, path)
					if err != nil {
						d.Logger.Errorf("cannot create target for path '%s': %v", path, err)
						return errors.Wrapf(err, "cannot create target for path '%s': %v", path, err)
					}
					defer func() {
						if err := targetFP.Close(); err != nil {
							d.Logger.Errorf("cannot close target: %v", err)
						}
					}()
					csWriter, err := checksum.NewChecksumWriter(
						[]checksum.DigestAlgorithm{checksum.DigestSHA512},
						targetFP,
					)
					if err != nil {
						d.Logger.Errorf("cannot create checksum writer for file '%s%s': %v", vfs, path, err)
						return errors.Wrapf(err, "cannot create checksum writer for file '%s%s': %v", vfs, path, err)
					}
					_, err = io.Copy(csWriter, sourceFP)
					if err != nil {
						if err := csWriter.Close(); err != nil {
							d.Logger.Errorf("cannot close checksum writer: %v", err)
							return errors.Wrapf(err, "cannot close checksum writer: %v", err)
						}
						d.Logger.Errorf("error writing file to path '%s%s': %v", vfs, path, err)
						return errors.Wrapf(err, "error writing file to path '%s%s': %v", vfs, path, err)
					}
					if err := csWriter.Close(); err != nil {
						d.Logger.Errorf("cannot close checksum writer: %v", err)
						return errors.Wrapf(err, "cannot close checksum writer: %v", err)
					}
					return nil
				}()
				if err != nil {
					d.Logger.Errorf("cannot copy object from: %s", pathToCopyFrom)
					continue
				}

				storagePartition.CurrentSize += int64(objectPb.Size)
				storagePartition.CurrentObjects++
				_, err = d.ClientStorageHandlerHandler.UpdateStoragePartition(ctx, storagePartition)
				if err != nil {
					d.Logger.Errorf("Could not update storage partition with ID: %v", storagePartition.Id)
					return &pb.NoParam{}, errors.Wrapf(err, "Could not update storage partition with ID: %v", storagePartition.Id)
				}
				objectInstance := &pb.ObjectInstance{Path: path, Status: "new", ObjectId: objectPb.Id, StoragePartitionId: storagePartition.Id, Size: objectPb.Size}
				_, err = d.ClientStorageHandlerHandler.CreateObjectInstance(ctx, objectInstance)
				if err != nil {
					d.Logger.Errorf("Could not create objectInstance for object with ID: %v", objectPb.Id)
					return &pb.NoParam{}, errors.Wrapf(err, "Could not create objectInstance for object with ID: %v", objectPb.Id)
				}
				//ToDo Decide what to do with deletion
				storageLocationsToDeleteFrom.StorageLocations = nil
				// Delete redundant objectInstances and objects from storageLocations
				if len(storageLocationsToDeleteFrom.StorageLocations) != 0 {
					for _, storageLocationToDeleteFrom := range storageLocationsToDeleteFrom.StorageLocations {
						storagePartitionsForLocationIdToDelete, err := d.ClientStorageHandlerHandler.GetStoragePartitionsByStorageLocationId(ctx, &pb.Id{Id: storageLocationToDeleteFrom.Id})
						if err != nil {
							d.Logger.Errorf("cannot get storagePartitions for storageLocation: %v", storageLocationToDeleteFrom.Alias)
							return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", storageLocationToDeleteFrom.Alias)
						}
						objectInstancesWithPartitionsToDelete := service.GetObjectInstancesToDelete(ObjectInstancesPb, storagePartitionsForLocationIdToDelete)
						for objectInstanceToDelete, partitionToDelete := range objectInstancesWithPartitionsToDelete {

							if err := writefs.Remove(vfs, objectInstanceToDelete.Path); err != nil {
								d.Logger.Errorf("error deleting file to '%s': %v", objectInstanceToDelete.Path, err)
								return &pb.NoParam{}, errors.Wrapf(err, "error writing file to '%s': %v", path, err)
							}

							_, err = d.ClientStorageHandlerHandler.DeleteObjectInstance(ctx, &pb.Id{Id: objectInstanceToDelete.Id})
							if err != nil {
								d.Logger.Errorf("Could not delete objectInstance with ID: %v", objectInstanceToDelete.Id)
								return &pb.NoParam{}, errors.Wrapf(err, "Could not delete objectInstance with ID: %v", objectInstanceToDelete.Id)
							}

							partitionToDelete.CurrentSize -= int64(objectInstanceToDelete.Size)
							partitionToDelete.CurrentObjects--
							_, err = d.ClientStorageHandlerHandler.UpdateStoragePartition(ctx, partitionToDelete)
							if err != nil {
								d.Logger.Errorf("Could not update storage partition with ID: %v", partitionToDelete.Id)
								return &pb.NoParam{}, errors.Wrapf(err, "Could not update storage partition with ID: %v", partitionToDelete.Id)
							}
						}
					}
				}
			}
			vfs.Close()
		}
	}

	return &pb.NoParam{}, nil
}
