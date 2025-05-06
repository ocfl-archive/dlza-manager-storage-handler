package main

import (
	archiveerror "github.com/ocfl-archive/error/pkg/error"
)

type errorID = archiveerror.ID

const (
	ErrorDataBaseRead           = "ErrorDataBaseRead"
	ErrorDataBaseWrite          = "ErrorDataBaseWrite"
	ErrorVfsOpen                = "ErrorVfsOpen"
	ErrorVfsClosingSource       = "ErrorVfsClosingSource"
	ErrorFileSystemCreateTarget = "ErrorFileSystemErrorFileSystemCreateTarget"
	ErrorFileSystemCloseTarget  = "ErrorFileSystemCloseTarget"
	ErrorIOCopy                 = "ErrorIOCopy"
)
