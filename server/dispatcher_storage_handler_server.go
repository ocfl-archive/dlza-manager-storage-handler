package server

import (
	"context"
	"encoding/json"
	"github.com/je4/utils/v2/pkg/checksum"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/models"
	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"io"
	"path/filepath"

	"github.com/je4/filesystem/v3/pkg/vfsrw"
	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/pkg/errors"
)

type DispatcherStorageHandlerServer struct {
	storageHandlerPb.UnimplementedDispatcherStorageHandlerServiceServer
	ClientStorageHandlerHandler handlerPb.StorageHandlerHandlerServiceClient
	Logger                      zLogger.ZLogger
}

func (d *DispatcherStorageHandlerServer) ChangeQualityForCollectionWithObjectIds(ctx context.Context, collectionsWithObjects *pb.CollectionAliases) (*pb.NoParam, error) {

	for _, collectionAlias := range collectionsWithObjects.CollectionAliases {
		//get cheapest storage locations
		storageLocationsPb, err := d.ClientStorageHandlerHandler.GetStorageLocationsByCollectionAlias(ctx, collectionAlias)
		if err != nil {
			return &pb.NoParam{}, errors.Wrapf(err, "cannot get storageLocations for collection: %v", collectionAlias)
		}
		if storageLocationsPb.StorageLocations == nil {
			d.Logger.Warn().Msgf("The collection " + collectionAlias.CollectionAlias + " doesn't have enough storage locations")
			continue
		}

		for _, objectId := range collectionAlias.Ids {

			objectPb, err := d.ClientStorageHandlerHandler.GetObjectById(ctx, objectId)
			if err != nil {
				return &pb.NoParam{}, errors.Wrapf(err, "cannot GetObjectById: %v", err)
			}

			currentStorageLocationsPb, err := d.ClientStorageHandlerHandler.GetStorageLocationsByObjectId(ctx, &pb.Id{Id: objectPb.Id})
			if err != nil {
				return &pb.NoParam{}, errors.Wrapf(err, "cannot GetCurrentStorageLocationsByCollectionAlias: %v", err)
			}

			storageLocationsToCopyIn := service.GetStorageLocationsToCopyIn(storageLocationsPb, currentStorageLocationsPb)
			storageLocationsToDeleteFrom := service.GetStorageLocationsToDeleteFrom(storageLocationsPb, currentStorageLocationsPb)
			if len(storageLocationsToCopyIn.StorageLocations) == 0 && len(storageLocationsToCopyIn.StorageLocations) == 0 {
				continue
			}
			configObj, err := models.LoadStorageHandlerConfig(append(currentStorageLocationsPb.StorageLocations, storageLocationsToCopyIn.StorageLocations...))
			if err != nil {
				return &pb.NoParam{}, errors.Wrapf(err, "cannot load storage-handler config: %v", err)
			}
			vfs, err := vfsrw.NewFS(configObj.VFS, d.Logger)
			if err != nil {
				d.Logger.Warn().Msgf("cannot create vfs: %v", err)
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
					d.Logger.Error().Msgf("cannot get storagePartition for storageLocation: %v", storageLocation.Alias)
					return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", storageLocation.Alias)
				}
				connection := models.Connection{}
				err = json.Unmarshal([]byte(storageLocation.Connection), &connection)
				if err != nil {
					d.Logger.Error().Msgf("error mapping json")
					return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", storageLocation.Alias)
				}

				path := connection.Folder + storagePartition.Alias + "/" + filepath.Base(pathToCopyFrom)
				sourceFP, err := vfs.Open(pathToCopyFrom)
				if err == nil {
					storagePartition.CurrentSize += int64(objectPb.Size)
					storagePartition.CurrentObjects++
					objectInstance := &pb.ObjectInstance{Path: path, Status: "new", ObjectId: objectPb.Id, StoragePartitionId: storagePartition.Id, Size: objectPb.Size}
					_, err = d.ClientStorageHandlerHandler.CreateObjectInstance(ctx, objectInstance)
					if err != nil {
						d.Logger.Error().Msgf("Could not create objectInstance for object with ID: %v", objectPb.Id)
						continue
					}
					_, err = d.ClientStorageHandlerHandler.UpdateStoragePartition(ctx, storagePartition)
					if err != nil {
						d.Logger.Error().Msgf("Could not update storage partition with ID: %v", storagePartition.Id)
						continue
					}
					err = func() error {

						targetFP, err := writefs.Create(vfs, path)
						if err != nil {
							d.Logger.Error().Msgf("cannot create target for path '%s': %v", path, err)
							sourceFP.Close()
							return errors.Wrapf(err, "cannot create target for path '%s': %v", path, err)
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
							d.Logger.Error().Msgf("cannot create checksum writer for file '%s%s': %v", vfs, path, err)
							sourceFP.Close()
							return errors.Wrapf(err, "cannot create checksum writer for file '%s%s': %v", vfs, path, err)
						}
						_, err = io.Copy(csWriter, sourceFP)
						if err != nil {
							if err := csWriter.Close(); err != nil {
								d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
								return errors.Wrapf(err, "cannot close checksum writer: %v", err)
							}
							d.Logger.Error().Msgf("error writing file to path '%s%s': %v", vfs, path, err)
							sourceFP.Close()
							return errors.Wrapf(err, "error writing file to path '%s%s': %v", vfs, path, err)
						}
						if err := csWriter.Close(); err != nil {
							d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
							return errors.Wrapf(err, "cannot close checksum writer: %v", err)
						}
						sourceFP.Close()
						return nil
					}()
					if err != nil {
						d.Logger.Error().Msgf("cannot copy object from: %s", pathToCopyFrom)
						continue
					}

					//ToDo Decide what to do with deletion
					storageLocationsToDeleteFrom.StorageLocations = nil
					// Delete redundant objectInstances and objects from storageLocations
					if len(storageLocationsToDeleteFrom.StorageLocations) != 0 {
						for _, storageLocationToDeleteFrom := range storageLocationsToDeleteFrom.StorageLocations {
							storagePartitionsForLocationIdToDelete, err := d.ClientStorageHandlerHandler.GetStoragePartitionsByStorageLocationId(ctx, &pb.Id{Id: storageLocationToDeleteFrom.Id})
							if err != nil {
								d.Logger.Error().Msgf("cannot get storagePartitions for storageLocation: %v", storageLocationToDeleteFrom.Alias)
								return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", storageLocationToDeleteFrom.Alias)
							}
							objectInstancesWithPartitionsToDelete := service.GetObjectInstancesToDelete(ObjectInstancesPb, storagePartitionsForLocationIdToDelete)
							for objectInstanceToDelete, partitionToDelete := range objectInstancesWithPartitionsToDelete {

								if err := writefs.Remove(vfs, objectInstanceToDelete.Path); err != nil {
									d.Logger.Error().Msgf("error deleting file to '%s': %v", objectInstanceToDelete.Path, err)
									return &pb.NoParam{}, errors.Wrapf(err, "error writing file to '%s': %v", path, err)
								}

								_, err = d.ClientStorageHandlerHandler.DeleteObjectInstance(ctx, &pb.Id{Id: objectInstanceToDelete.Id})
								if err != nil {
									d.Logger.Error().Msgf("Could not delete objectInstance with ID: %v", objectInstanceToDelete.Id)
									return &pb.NoParam{}, errors.Wrapf(err, "Could not delete objectInstance with ID: %v", objectInstanceToDelete.Id)
								}

								partitionToDelete.CurrentSize -= int64(objectInstanceToDelete.Size)
								partitionToDelete.CurrentObjects--
								_, err = d.ClientStorageHandlerHandler.UpdateStoragePartition(ctx, partitionToDelete)
								if err != nil {
									d.Logger.Error().Msgf("Could not update storage partition with ID: %v", partitionToDelete.Id)
									return &pb.NoParam{}, errors.Wrapf(err, "Could not update storage partition with ID: %v", partitionToDelete.Id)
								}
							}
						}
					}
				} else {
					d.Logger.Error().Msgf("cannot open file '%s': %v", pathToCopyFrom, err)
					vfs.Close()
					continue
				}
			}
			vfs.Close()
		}
	}

	return &pb.NoParam{}, nil
}

func (d *DispatcherStorageHandlerServer) CopyArchiveTo(ctx context.Context, copyFromTo *pb.CopyFromTo) (*pb.NoParam, error) {
	storagePartition, err := d.ClientStorageHandlerHandler.GetStoragePartitionForLocation(ctx, &pb.SizeAndId{Size: copyFromTo.ObjectInstance.Size, Id: copyFromTo.LocationCopyTo.Id})
	if err != nil {
		d.Logger.Error().Msgf("cannot get storagePartition for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
		return &pb.NoParam{}, errors.Wrapf(err, "cannot get storagePartition for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	storageLocationToCopyFrom, err := d.ClientStorageHandlerHandler.GetStorageLocationByObjectInstanceId(ctx, &pb.Id{Id: copyFromTo.ObjectInstance.Id})
	if err != nil {
		d.Logger.Error().Msgf("cannot get GetStorageLocationByObjectInstanceId for object instance ID: %v", copyFromTo.ObjectInstance.Id)
		return &pb.NoParam{}, errors.Wrapf(err, "cannot get GetStorageLocationByObjectInstanceId for object instance ID: %v", copyFromTo.ObjectInstance.Id)
	}

	connection := models.Connection{}
	err = json.Unmarshal([]byte(copyFromTo.LocationCopyTo.Connection), &connection)
	if err != nil {
		d.Logger.Error().Msgf("error mapping json")
		return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	vfsConfig, err := config.LoadVfsConfig(copyFromTo.LocationCopyTo.Connection, storageLocationToCopyFrom.Connection)
	if err != nil {
		d.Logger.Error().Msgf("error mapping json for storage location connection field: %v", err)
		return nil, errors.Wrapf(err, "error mapping json for storage location connection field")
	}
	vfs, err := vfsrw.NewFS(vfsConfig, d.Logger)
	if err != nil {
		d.Logger.Warn().Msgf("cannot create vfs: %v", err)
		return &pb.NoParam{}, errors.Wrapf(err, "error mapping json for storageLocation: %v", copyFromTo.LocationCopyTo.Alias)
	}

	path := connection.Folder + storagePartition.Alias + "/" + filepath.Base(copyFromTo.ObjectInstance.Path)
	sourceFP, err := vfs.Open(copyFromTo.ObjectInstance.Path)

	if err == nil {
		storagePartition.CurrentSize += copyFromTo.ObjectInstance.Size
		storagePartition.CurrentObjects++
		objectInstance := &pb.ObjectInstance{Path: path, Status: "new", ObjectId: copyFromTo.ObjectInstance.ObjectId, StoragePartitionId: storagePartition.Id, Size: copyFromTo.ObjectInstance.Size}
		_, err = d.ClientStorageHandlerHandler.CreateObjectInstance(ctx, objectInstance)
		if err != nil {
			d.Logger.Error().Msgf("Could not create objectInstance for object with ID: %v", copyFromTo.ObjectInstance.ObjectId)
			return &pb.NoParam{}, errors.Wrapf(err, "Could not create objectInstance for object with ID: %v", copyFromTo.ObjectInstance.ObjectId)
		}
		_, err = d.ClientStorageHandlerHandler.UpdateStoragePartition(ctx, storagePartition)
		if err != nil {
			d.Logger.Error().Msgf("Could not update storage partition with ID: %v", storagePartition.Id)
			return &pb.NoParam{}, errors.Wrapf(err, "Could not update storage partition with ID: %v", storagePartition.Id)
		}

		targetFP, err := writefs.Create(vfs, path)
		if err != nil {
			d.Logger.Error().Msgf("cannot create target for path '%v': %s", path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "cannot create target for path '%v': %s", path, err)
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
			d.Logger.Error().Msgf("cannot create checksum writer for file '%s%s': %v", vfs, path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "cannot create checksum writer for file '%v%v': %s", vfs, path, err)
		}
		_, err = io.Copy(csWriter, sourceFP)
		if err != nil {
			if err := csWriter.Close(); err != nil {
				d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
				return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
			}
			d.Logger.Error().Msgf("error writing file to path '%v%v': %s", vfs, path, err)
			sourceFP.Close()
			return &pb.NoParam{}, errors.Wrapf(err, "error writing file to path '%v%v': %s", vfs, path, err)
		}
		if err := csWriter.Close(); err != nil {
			d.Logger.Error().Msgf("cannot close checksum writer: %v", err)
			return &pb.NoParam{}, errors.Wrapf(err, "cannot close checksum writer: %v", err)
		}
		sourceFP.Close()
		vfs.Close()
	} else {
		d.Logger.Error().Msgf("error opening file with path %v: %s", copyFromTo.ObjectInstance.Path, err)
		return nil, errors.Wrapf(err, "error opening file with path %v", copyFromTo.ObjectInstance.Path)
	}
	return &pb.NoParam{}, nil
}
