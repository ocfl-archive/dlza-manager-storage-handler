package models

import (
	"database/sql"
	"github.com/je4/filesystem/v3/pkg/vfsrw"
)

const (
	sftp      string = "sftp"
	local     string = "local"
	s3_switch string = "s3_switch"
	s3_amazon string = "s3_amazon"
)

type Connection struct {
	Folder string
	VFS    vfsrw.Config
}

type StorageLocation struct {
	Alias              string
	Type               string
	Vault              sql.NullString
	Connection         string
	Quality            int
	Price              int
	SecurityCompliency string
	FillFirst          bool
	OcflType           string
	TenantId           string
	Id                 string
	NumberOfThreads    int
	// virtual values
	TotalFilesSize      int64
	TotalExistingVolume int64
}
