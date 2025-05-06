package service

import (
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"testing"
)

func TestGetStorageLocationsToCopyIn(t *testing.T) {

	storageLocationsNeeded := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}, {Type: "Amazon S3", Id: "Amazon S3"}, {Type: "Switch S3", Id: "Switch S3"}, {Type: "local", Id: "local"}}}
	currentStorageLocations := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}}}

	storageLocationsToCopyIn := GetStorageLocationsToCopyIn(storageLocationsNeeded, currentStorageLocations)
	if len(storageLocationsToCopyIn.StorageLocations) != 3 {
		panic("TestGetStorageLocationsToCopyIn has failed")
	}
	for _, storageLocation := range storageLocationsToCopyIn.StorageLocations {
		if storageLocation.Type == "sftp" {
			panic("TestGetStorageLocationsToCopyIn has failed")
		}
	}
}

func TestGetStorageLocationsToCopyIn2(t *testing.T) {

	storageLocationsNeeded := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}, {Type: "Amazon S3", Id: "Amazon S3"}, {Type: "Switch S3", Id: "Switch S3"}, {Type: "local", Id: "local"}}}
	currentStorageLocations := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}, {Type: "Amazon S3", Id: "Amazon S3"}}}

	storageLocationsToCopyIn := GetStorageLocationsToCopyIn(storageLocationsNeeded, currentStorageLocations)
	if len(storageLocationsToCopyIn.StorageLocations) != 2 {
		panic("TestGetStorageLocationsToCopyIn2 has failed")
	}
	for _, storageLocation := range storageLocationsToCopyIn.StorageLocations {
		if storageLocation.Type == "sftp" || storageLocation.Type == "Amazon S3" {
			panic("TestGetStorageLocationsToCopyIn2 has failed")
		}
	}
}

func TestGetStorageLocationsToCopyIn3(t *testing.T) {

	storageLocationsNeeded := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}}}
	currentStorageLocations := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "Amazon S3", Id: "Amazon S3"}, {Type: "Switch S3", Id: "Switch S3"}, {Type: "local", Id: "local"}}}

	storageLocationsToCopyIn := GetStorageLocationsToCopyIn(storageLocationsNeeded, currentStorageLocations)
	if len(storageLocationsToCopyIn.StorageLocations) != 1 {
		panic("TestGetStorageLocationsToCopyIn3 has failed")
	}
	for _, storageLocation := range storageLocationsToCopyIn.StorageLocations {
		if storageLocation.Type != "sftp" {
			panic("TestGetStorageLocationsToCopyIn3 has failed")
		}
	}
}

func TestGetStorageLocationsToDeleteFrom(t *testing.T) {

	storageLocationsNeeded := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}, {Type: "Amazon S3", Id: "Amazon S3"}, {Type: "Switch S3", Id: "Switch S3"}, {Type: "local", Id: "local"}}}
	currentStorageLocations := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}}}

	storageLocationsToDeleteFrom := GetStorageLocationsToDeleteFrom(storageLocationsNeeded, currentStorageLocations)
	if len(storageLocationsToDeleteFrom.StorageLocations) != 0 {
		panic("TestGetStorageLocationsToDeleteFrom has failed")
	}
}

func TestGetStorageLocationsToDeleteFrom2(t *testing.T) {

	storageLocationsNeeded := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "Switch S3", Id: "Switch S3"}, {Type: "local", Id: "local"}}}
	currentStorageLocations := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}, {Type: "Amazon S3", Id: "Amazon S3"}}}

	storageLocationsToDeleteFrom := GetStorageLocationsToDeleteFrom(storageLocationsNeeded, currentStorageLocations)
	if len(storageLocationsToDeleteFrom.StorageLocations) != 2 {
		panic("TestGetStorageLocationsToDeleteFrom2 has failed")
	}
	for _, storageLocation := range storageLocationsToDeleteFrom.StorageLocations {
		if storageLocation.Type == "Switch S3" || storageLocation.Type == "local" {
			panic("TestGetStorageLocationsToDeleteFrom2 has failed")
		}
	}
}

func TestGetStorageLocationsToDeleteFrom3(t *testing.T) {

	storageLocationsNeeded := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "sftp", Id: "sftp"}}}
	currentStorageLocations := &pb.StorageLocations{StorageLocations: []*pb.StorageLocation{{Type: "Amazon S3", Id: "Amazon S3"}, {Type: "Switch S3", Id: "Switch S3"}, {Type: "local", Id: "local"}}}

	storageLocationsToCopyIn := GetStorageLocationsToDeleteFrom(storageLocationsNeeded, currentStorageLocations)
	if len(storageLocationsToCopyIn.StorageLocations) != 3 {
		panic("TestGetStorageLocationsToDeleteFrom3 has failed")
	}
	for _, storageLocation := range storageLocationsToCopyIn.StorageLocations {
		if storageLocation.Type == "sftp" {
			panic("TestGetStorageLocationsToDeleteFrom3 has failed")
		}
	}
}

func TestGetObjectInstancesToDelete(t *testing.T) {

	storagePartitions := &pb.StoragePartitions{StoragePartitions: []*pb.StoragePartition{{Id: "1"}}}
	objectInstances := &pb.ObjectInstances{ObjectInstances: []*pb.ObjectInstance{{StoragePartitionId: "1"}, {StoragePartitionId: "2"}, {StoragePartitionId: "3"}}}

	objectInstancesWithPartitionToDelete := GetObjectInstancesToDelete(objectInstances, storagePartitions)
	if len(objectInstancesWithPartitionToDelete) != 1 {
		panic("TestGetObjectInstancesToDelete has failed")
	}
	for objectInstance, partition := range objectInstancesWithPartitionToDelete {
		if objectInstance.StoragePartitionId != partition.Id && objectInstance.StoragePartitionId != "1" {
			panic("TestGetObjectInstancesToDelete has failed")
		}
	}
}

func TestGetObjectInstancesToDelete2(t *testing.T) {

	storagePartitions := &pb.StoragePartitions{StoragePartitions: []*pb.StoragePartition{{Id: "1"}, {Id: "2"}, {Id: "3"}}}
	objectInstances := &pb.ObjectInstances{ObjectInstances: []*pb.ObjectInstance{{StoragePartitionId: "1"}, {StoragePartitionId: "2"}, {StoragePartitionId: "3"}}}

	objectInstancesWithPartitionToDelete := GetObjectInstancesToDelete(objectInstances, storagePartitions)
	if len(objectInstancesWithPartitionToDelete) != 3 {
		panic("TestGetObjectInstancesToDelete2 has failed")
	}
	for objectInstance, partition := range objectInstancesWithPartitionToDelete {
		if objectInstance.StoragePartitionId != partition.Id {
			panic("TestGetObjectInstancesToDelete2 has failed")
		}
		if objectInstance.StoragePartitionId != "1" && objectInstance.StoragePartitionId != "2" && objectInstance.StoragePartitionId != "3" {
			panic("TestGetObjectInstancesToDelete2 has failed")
		}
	}
}

func TestGetObjectInstancesToDelete3(t *testing.T) {

	storagePartitions := &pb.StoragePartitions{StoragePartitions: []*pb.StoragePartition{{Id: "1"}, {Id: "5"}}}
	objectInstances := &pb.ObjectInstances{ObjectInstances: []*pb.ObjectInstance{{StoragePartitionId: "1"}, {StoragePartitionId: "2"}, {StoragePartitionId: "3"}, {StoragePartitionId: "4"}, {StoragePartitionId: "5"}}}

	objectInstancesWithPartitionToDelete := GetObjectInstancesToDelete(objectInstances, storagePartitions)
	if len(objectInstancesWithPartitionToDelete) != 2 {
		panic("TestGetObjectInstancesToDelete3 has failed")
	}
	for objectInstance, partition := range objectInstancesWithPartitionToDelete {
		if objectInstance.StoragePartitionId != partition.Id {
			panic("TestGetObjectInstancesToDelete3 has failed")
		}
		if objectInstance.StoragePartitionId != "1" && objectInstance.StoragePartitionId != "5" {
			panic("TestGetObjectInstancesToDelete3 has failed")
		}
	}
}
