// protoc --go_out=. --go-grpc_out=. proto/copy.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: storage_handler_proto.proto

package storagehandlerproto

import (
	dlzamanagerproto "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_storage_handler_proto_proto protoreflect.FileDescriptor

var file_storage_handler_proto_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x10, 0x64, 0x6c, 0x7a, 0x61, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0x68, 0x0a, 0x1d, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x08, 0x43, 0x6f, 0x70, 0x79, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x1f, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x32, 0x88,
	0x02, 0x0a, 0x16, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x0f, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x26, 0x2e, 0x64,
	0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00,
	0x12, 0x4c, 0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x41, 0x6e,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x41, 0x6e, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x12, 0x49,
	0x0a, 0x0b, 0x41, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x1a, 0x18, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x32, 0x76, 0x0a, 0x1a, 0x43, 0x6c, 0x65,
	0x72, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x58, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x22, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x74,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x00, 0x32, 0xef, 0x16, 0x0a, 0x13, 0x43, 0x6c, 0x65, 0x72, 0x6b, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x0e, 0x46, 0x69, 0x6e,
	0x64, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49,
	0x64, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x40, 0x0a,
	0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x14, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x49, 0x64, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12,
	0x42, 0x0a, 0x0a, 0x53, 0x61, 0x76, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x18, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x12, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x1a, 0x18, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0e, 0x46, 0x69, 0x6e,
	0x64, 0x41, 0x6c, 0x6c, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x19, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e,
	0x6f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a, 0x19, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x73, 0x22, 0x00, 0x12, 0x5b, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x22, 0x2e, 0x64, 0x6c, 0x7a,
	0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x00,
	0x12, 0x54, 0x0a, 0x13, 0x53, 0x61, 0x76, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a,
	0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x1d, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x49, 0x64, 0x1a, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49,
	0x64, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x4c, 0x0a,
	0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x10, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x18, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0d, 0x47, 0x65, 0x74,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a,
	0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64,
	0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x15,
	0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x20, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x00, 0x12,
	0x3d, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14,
	0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x16, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x00, 0x12, 0x5b,
	0x0a, 0x1a, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64,
	0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x49, 0x64, 0x1a, 0x25, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x21, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00,
	0x12, 0x55, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x49, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49,
	0x64, 0x1a, 0x22, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x74,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x17, 0x46, 0x69, 0x6e, 0x64, 0x41,
	0x6c, 0x6c, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x19, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x22, 0x00, 0x12, 0x62, 0x0a,
	0x21, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42,
	0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x1d, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x00, 0x12, 0x5e, 0x0a, 0x21, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x42,
	0x79, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x19, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x22,
	0x00, 0x12, 0x5a, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x50, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x17, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x00, 0x12, 0x5a, 0x0a,
	0x1b, 0x47, 0x65, 0x74, 0x4d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x46, 0x6f, 0x72,
	0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x64,
	0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1b, 0x2e, 0x64, 0x6c, 0x7a,
	0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69,
	0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x19, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x6f, 0x6e, 0x6f, 0x6d, 0x73, 0x46, 0x6f, 0x72, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x19, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x6e, 0x6f, 0x6d, 0x73, 0x22,
	0x00, 0x12, 0x6a, 0x0a, 0x25, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x42, 0x79, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49,
	0x64, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a,
	0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x21, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x22, 0x00, 0x12, 0x56, 0x0a,
	0x1b, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64,
	0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x17, 0x2e, 0x64, 0x6c, 0x7a,
	0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x22, 0x00, 0x12, 0x7c, 0x0a, 0x32, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x73,
	0x42, 0x79, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x49, 0x64, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x26, 0x2e, 0x64, 0x6c, 0x7a, 0x61,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x73, 0x22, 0x00, 0x12, 0x6c, 0x0a, 0x26, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x49, 0x64, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e,
	0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x22, 0x2e, 0x64, 0x6c,
	0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x00, 0x12, 0x70, 0x0a, 0x29, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50,
	0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c,
	0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x23, 0x2e, 0x64,
	0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0x00, 0x12, 0x74, 0x0a, 0x2f, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x42, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1c, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x21, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0b, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x1a, 0x1e,
	0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x00,
	0x12, 0x46, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x1e, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x1a, 0x14, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x0b, 0x41, 0x6c, 0x74, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x32, 0x8e, 0x01, 0x0a, 0x1f, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6b, 0x0a, 0x27, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x46, 0x6f, 0x72, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x74, 0x68, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49,
	0x64, 0x73, 0x12, 0x23, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x41, 0x6c, 0x69, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x19, 0x2e, 0x64, 0x6c, 0x7a, 0x61, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x22, 0x00, 0x42, 0xa5, 0x01, 0x0a, 0x1e, 0x63, 0x68, 0x2e, 0x75, 0x6e, 0x69, 0x62,
	0x61, 0x73, 0x2e, 0x75, 0x62, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x68, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x67, 0x42, 0x13, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x48,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x63, 0x66, 0x6c, 0x2d,
	0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x2f, 0x64, 0x6c, 0x7a, 0x61, 0x2d, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2d, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2d, 0x68, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x68, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0xa2, 0x02, 0x03, 0x55, 0x42, 0x42, 0xaa, 0x02,
	0x1b, 0x55, 0x6e, 0x69, 0x62, 0x61, 0x73, 0x2e, 0x55, 0x42, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x2e, 0x50, 0x47, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_storage_handler_proto_proto_goTypes = []interface{}{
	(*dlzamanagerproto.IncomingOrder)(nil),        // 0: dlzamanagerproto.IncomingOrder
	(*dlzamanagerproto.UploaderAccessObject)(nil), // 1: dlzamanagerproto.UploaderAccessObject
	(*dlzamanagerproto.ObjectAndFiles)(nil),       // 2: dlzamanagerproto.ObjectAndFiles
	(*dlzamanagerproto.StatusObject)(nil),         // 3: dlzamanagerproto.StatusObject
	(*dlzamanagerproto.StoragePartition)(nil),     // 4: dlzamanagerproto.StoragePartition
	(*dlzamanagerproto.Id)(nil),                   // 5: dlzamanagerproto.Id
	(*dlzamanagerproto.Tenant)(nil),               // 6: dlzamanagerproto.Tenant
	(*dlzamanagerproto.NoParam)(nil),              // 7: dlzamanagerproto.NoParam
	(*dlzamanagerproto.StorageLocation)(nil),      // 8: dlzamanagerproto.StorageLocation
	(*dlzamanagerproto.Collection)(nil),           // 9: dlzamanagerproto.Collection
	(*dlzamanagerproto.Pagination)(nil),           // 10: dlzamanagerproto.Pagination
	(*dlzamanagerproto.CollectionAliases)(nil),    // 11: dlzamanagerproto.CollectionAliases
	(*dlzamanagerproto.Status)(nil),               // 12: dlzamanagerproto.Status
	(*dlzamanagerproto.Tenants)(nil),              // 13: dlzamanagerproto.Tenants
	(*dlzamanagerproto.StorageLocations)(nil),     // 14: dlzamanagerproto.StorageLocations
	(*dlzamanagerproto.Collections)(nil),          // 15: dlzamanagerproto.Collections
	(*dlzamanagerproto.Object)(nil),               // 16: dlzamanagerproto.Object
	(*dlzamanagerproto.ObjectInstance)(nil),       // 17: dlzamanagerproto.ObjectInstance
	(*dlzamanagerproto.File)(nil),                 // 18: dlzamanagerproto.File
	(*dlzamanagerproto.ObjectInstanceCheck)(nil),  // 19: dlzamanagerproto.ObjectInstanceCheck
	(*dlzamanagerproto.Objects)(nil),              // 20: dlzamanagerproto.Objects
	(*dlzamanagerproto.Files)(nil),                // 21: dlzamanagerproto.Files
	(*dlzamanagerproto.MimeTypes)(nil),            // 22: dlzamanagerproto.MimeTypes
	(*dlzamanagerproto.Pronoms)(nil),              // 23: dlzamanagerproto.Pronoms
	(*dlzamanagerproto.ObjectInstances)(nil),      // 24: dlzamanagerproto.ObjectInstances
	(*dlzamanagerproto.ObjectInstanceChecks)(nil), // 25: dlzamanagerproto.ObjectInstanceChecks
	(*dlzamanagerproto.StoragePartitions)(nil),    // 26: dlzamanagerproto.StoragePartitions
}
var file_storage_handler_proto_proto_depIdxs = []int32{
	0,  // 0: storagehandlerproto.UploaderStorageHandlerService.CopyFile:input_type -> dlzamanagerproto.IncomingOrder
	1,  // 1: storagehandlerproto.UploaderHandlerService.TenantHasAccess:input_type -> dlzamanagerproto.UploaderAccessObject
	2,  // 2: storagehandlerproto.UploaderHandlerService.SaveObjectAndFiles:input_type -> dlzamanagerproto.ObjectAndFiles
	3,  // 3: storagehandlerproto.UploaderHandlerService.AlterStatus:input_type -> dlzamanagerproto.StatusObject
	4,  // 4: storagehandlerproto.ClerkStorageHandlerService.CreateStoragePartition:input_type -> dlzamanagerproto.StoragePartition
	5,  // 5: storagehandlerproto.ClerkHandlerService.FindTenantById:input_type -> dlzamanagerproto.Id
	5,  // 6: storagehandlerproto.ClerkHandlerService.DeleteTenant:input_type -> dlzamanagerproto.Id
	6,  // 7: storagehandlerproto.ClerkHandlerService.SaveTenant:input_type -> dlzamanagerproto.Tenant
	6,  // 8: storagehandlerproto.ClerkHandlerService.UpdateTenant:input_type -> dlzamanagerproto.Tenant
	7,  // 9: storagehandlerproto.ClerkHandlerService.FindAllTenants:input_type -> dlzamanagerproto.NoParam
	5,  // 10: storagehandlerproto.ClerkHandlerService.GetStorageLocationsByTenantId:input_type -> dlzamanagerproto.Id
	8,  // 11: storagehandlerproto.ClerkHandlerService.SaveStorageLocation:input_type -> dlzamanagerproto.StorageLocation
	5,  // 12: storagehandlerproto.ClerkHandlerService.DeleteStorageLocationById:input_type -> dlzamanagerproto.Id
	5,  // 13: storagehandlerproto.ClerkHandlerService.GetCollectionsByTenantId:input_type -> dlzamanagerproto.Id
	5,  // 14: storagehandlerproto.ClerkHandlerService.GetCollectionById:input_type -> dlzamanagerproto.Id
	5,  // 15: storagehandlerproto.ClerkHandlerService.DeleteCollectionById:input_type -> dlzamanagerproto.Id
	9,  // 16: storagehandlerproto.ClerkHandlerService.CreateCollection:input_type -> dlzamanagerproto.Collection
	9,  // 17: storagehandlerproto.ClerkHandlerService.UpdateCollection:input_type -> dlzamanagerproto.Collection
	5,  // 18: storagehandlerproto.ClerkHandlerService.GetObjectById:input_type -> dlzamanagerproto.Id
	5,  // 19: storagehandlerproto.ClerkHandlerService.GetObjectInstanceById:input_type -> dlzamanagerproto.Id
	5,  // 20: storagehandlerproto.ClerkHandlerService.GetFileById:input_type -> dlzamanagerproto.Id
	5,  // 21: storagehandlerproto.ClerkHandlerService.GetObjectInstanceCheckById:input_type -> dlzamanagerproto.Id
	5,  // 22: storagehandlerproto.ClerkHandlerService.GetStorageLocationById:input_type -> dlzamanagerproto.Id
	5,  // 23: storagehandlerproto.ClerkHandlerService.GetStoragePartitionById:input_type -> dlzamanagerproto.Id
	10, // 24: storagehandlerproto.ClerkHandlerService.FindAllTenantsPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 25: storagehandlerproto.ClerkHandlerService.GetCollectionsByTenantIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 26: storagehandlerproto.ClerkHandlerService.GetObjectsByCollectionIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 27: storagehandlerproto.ClerkHandlerService.GetFilesByCollectionIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 28: storagehandlerproto.ClerkHandlerService.GetMimeTypesForCollectionId:input_type -> dlzamanagerproto.Pagination
	10, // 29: storagehandlerproto.ClerkHandlerService.GetPronomsForCollectionId:input_type -> dlzamanagerproto.Pagination
	10, // 30: storagehandlerproto.ClerkHandlerService.GetObjectInstancesByObjectIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 31: storagehandlerproto.ClerkHandlerService.GetFilesByObjectIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 32: storagehandlerproto.ClerkHandlerService.GetObjectInstanceChecksByObjectInstanceIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 33: storagehandlerproto.ClerkHandlerService.GetStorageLocationsByTenantIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 34: storagehandlerproto.ClerkHandlerService.GetStoragePartitionsByLocationIdPaginated:input_type -> dlzamanagerproto.Pagination
	10, // 35: storagehandlerproto.ClerkHandlerService.GetObjectInstancesByStoragePartitionIdPaginated:input_type -> dlzamanagerproto.Pagination
	5,  // 36: storagehandlerproto.ClerkHandlerService.CheckStatus:input_type -> dlzamanagerproto.Id
	3,  // 37: storagehandlerproto.ClerkHandlerService.CreateStatus:input_type -> dlzamanagerproto.StatusObject
	3,  // 38: storagehandlerproto.ClerkHandlerService.AlterStatus:input_type -> dlzamanagerproto.StatusObject
	11, // 39: storagehandlerproto.DispatcherStorageHandlerService.ChangeQualityForCollectionWithObjectIds:input_type -> dlzamanagerproto.CollectionAliases
	12, // 40: storagehandlerproto.UploaderStorageHandlerService.CopyFile:output_type -> dlzamanagerproto.Status
	12, // 41: storagehandlerproto.UploaderHandlerService.TenantHasAccess:output_type -> dlzamanagerproto.Status
	5,  // 42: storagehandlerproto.UploaderHandlerService.SaveObjectAndFiles:output_type -> dlzamanagerproto.Id
	12, // 43: storagehandlerproto.UploaderHandlerService.AlterStatus:output_type -> dlzamanagerproto.Status
	12, // 44: storagehandlerproto.ClerkStorageHandlerService.CreateStoragePartition:output_type -> dlzamanagerproto.Status
	6,  // 45: storagehandlerproto.ClerkHandlerService.FindTenantById:output_type -> dlzamanagerproto.Tenant
	12, // 46: storagehandlerproto.ClerkHandlerService.DeleteTenant:output_type -> dlzamanagerproto.Status
	12, // 47: storagehandlerproto.ClerkHandlerService.SaveTenant:output_type -> dlzamanagerproto.Status
	12, // 48: storagehandlerproto.ClerkHandlerService.UpdateTenant:output_type -> dlzamanagerproto.Status
	13, // 49: storagehandlerproto.ClerkHandlerService.FindAllTenants:output_type -> dlzamanagerproto.Tenants
	14, // 50: storagehandlerproto.ClerkHandlerService.GetStorageLocationsByTenantId:output_type -> dlzamanagerproto.StorageLocations
	12, // 51: storagehandlerproto.ClerkHandlerService.SaveStorageLocation:output_type -> dlzamanagerproto.Status
	12, // 52: storagehandlerproto.ClerkHandlerService.DeleteStorageLocationById:output_type -> dlzamanagerproto.Status
	15, // 53: storagehandlerproto.ClerkHandlerService.GetCollectionsByTenantId:output_type -> dlzamanagerproto.Collections
	9,  // 54: storagehandlerproto.ClerkHandlerService.GetCollectionById:output_type -> dlzamanagerproto.Collection
	12, // 55: storagehandlerproto.ClerkHandlerService.DeleteCollectionById:output_type -> dlzamanagerproto.Status
	12, // 56: storagehandlerproto.ClerkHandlerService.CreateCollection:output_type -> dlzamanagerproto.Status
	12, // 57: storagehandlerproto.ClerkHandlerService.UpdateCollection:output_type -> dlzamanagerproto.Status
	16, // 58: storagehandlerproto.ClerkHandlerService.GetObjectById:output_type -> dlzamanagerproto.Object
	17, // 59: storagehandlerproto.ClerkHandlerService.GetObjectInstanceById:output_type -> dlzamanagerproto.ObjectInstance
	18, // 60: storagehandlerproto.ClerkHandlerService.GetFileById:output_type -> dlzamanagerproto.File
	19, // 61: storagehandlerproto.ClerkHandlerService.GetObjectInstanceCheckById:output_type -> dlzamanagerproto.ObjectInstanceCheck
	8,  // 62: storagehandlerproto.ClerkHandlerService.GetStorageLocationById:output_type -> dlzamanagerproto.StorageLocation
	4,  // 63: storagehandlerproto.ClerkHandlerService.GetStoragePartitionById:output_type -> dlzamanagerproto.StoragePartition
	13, // 64: storagehandlerproto.ClerkHandlerService.FindAllTenantsPaginated:output_type -> dlzamanagerproto.Tenants
	15, // 65: storagehandlerproto.ClerkHandlerService.GetCollectionsByTenantIdPaginated:output_type -> dlzamanagerproto.Collections
	20, // 66: storagehandlerproto.ClerkHandlerService.GetObjectsByCollectionIdPaginated:output_type -> dlzamanagerproto.Objects
	21, // 67: storagehandlerproto.ClerkHandlerService.GetFilesByCollectionIdPaginated:output_type -> dlzamanagerproto.Files
	22, // 68: storagehandlerproto.ClerkHandlerService.GetMimeTypesForCollectionId:output_type -> dlzamanagerproto.MimeTypes
	23, // 69: storagehandlerproto.ClerkHandlerService.GetPronomsForCollectionId:output_type -> dlzamanagerproto.Pronoms
	24, // 70: storagehandlerproto.ClerkHandlerService.GetObjectInstancesByObjectIdPaginated:output_type -> dlzamanagerproto.ObjectInstances
	21, // 71: storagehandlerproto.ClerkHandlerService.GetFilesByObjectIdPaginated:output_type -> dlzamanagerproto.Files
	25, // 72: storagehandlerproto.ClerkHandlerService.GetObjectInstanceChecksByObjectInstanceIdPaginated:output_type -> dlzamanagerproto.ObjectInstanceChecks
	14, // 73: storagehandlerproto.ClerkHandlerService.GetStorageLocationsByTenantIdPaginated:output_type -> dlzamanagerproto.StorageLocations
	26, // 74: storagehandlerproto.ClerkHandlerService.GetStoragePartitionsByLocationIdPaginated:output_type -> dlzamanagerproto.StoragePartitions
	24, // 75: storagehandlerproto.ClerkHandlerService.GetObjectInstancesByStoragePartitionIdPaginated:output_type -> dlzamanagerproto.ObjectInstances
	3,  // 76: storagehandlerproto.ClerkHandlerService.CheckStatus:output_type -> dlzamanagerproto.StatusObject
	5,  // 77: storagehandlerproto.ClerkHandlerService.CreateStatus:output_type -> dlzamanagerproto.Id
	12, // 78: storagehandlerproto.ClerkHandlerService.AlterStatus:output_type -> dlzamanagerproto.Status
	7,  // 79: storagehandlerproto.DispatcherStorageHandlerService.ChangeQualityForCollectionWithObjectIds:output_type -> dlzamanagerproto.NoParam
	40, // [40:80] is the sub-list for method output_type
	0,  // [0:40] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_storage_handler_proto_proto_init() }
func file_storage_handler_proto_proto_init() {
	if File_storage_handler_proto_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_storage_handler_proto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   5,
		},
		GoTypes:           file_storage_handler_proto_proto_goTypes,
		DependencyIndexes: file_storage_handler_proto_proto_depIdxs,
	}.Build()
	File_storage_handler_proto_proto = out.File
	file_storage_handler_proto_proto_rawDesc = nil
	file_storage_handler_proto_proto_goTypes = nil
	file_storage_handler_proto_proto_depIdxs = nil
}
