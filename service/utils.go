package service

import (
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
)

func GetStorageLocationsToCopyIn(storageLocationsNeeded *pb.StorageLocations, currentStorageLocations *pb.StorageLocations) *pb.StorageLocations {

	resultingStorageLocations := make([]*pb.StorageLocation, 0)

	for _, storageLocationNeeded := range storageLocationsNeeded.StorageLocations {
		for index, currentStorageLocation := range currentStorageLocations.StorageLocations {
			if storageLocationNeeded.Id == currentStorageLocation.Id {
				break
			}
			if index == len(currentStorageLocations.StorageLocations)-1 {
				resultingStorageLocations = append(resultingStorageLocations, storageLocationNeeded)
			}
		}
	}
	return &pb.StorageLocations{StorageLocations: resultingStorageLocations}
}

func GetStorageLocationsToDeleteFrom(storageLocationsNeeded *pb.StorageLocations, currentStorageLocations *pb.StorageLocations) *pb.StorageLocations {

	storageLocationsToClean := make([]*pb.StorageLocation, 0)
	for _, currentStorageLocation := range currentStorageLocations.StorageLocations {
		for index, storageLocationNeeded := range storageLocationsNeeded.StorageLocations {
			if currentStorageLocation.Id == storageLocationNeeded.Id {
				break
			}
			if index == len(storageLocationsNeeded.StorageLocations)-1 {
				storageLocationsToClean = append(storageLocationsToClean, currentStorageLocation)
			}
		}
	}
	return &pb.StorageLocations{StorageLocations: storageLocationsToClean}
}

func GetObjectInstancesToDelete(objectInstances *pb.ObjectInstances, storagePartitionsToDeleteFrom *pb.StoragePartitions) map[*pb.ObjectInstance]*pb.StoragePartition {
	objectInstancesWithPartitionToDelete := make(map[*pb.ObjectInstance]*pb.StoragePartition)
	for _, objectInstance := range objectInstances.ObjectInstances {
		for _, storagePartition := range storagePartitionsToDeleteFrom.StoragePartitions {
			if storagePartition.Id == objectInstance.StoragePartitionId {
				objectInstancesWithPartitionToDelete[objectInstance] = storagePartition
				break
			}
		}
	}
	return objectInstancesWithPartitionToDelete
}
