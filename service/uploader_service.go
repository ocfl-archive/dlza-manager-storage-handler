package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"time"

	"emperror.dev/errors"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/dlza-manager/mapper"
	"github.com/ocfl-archive/dlza-manager/models"
	"github.com/ocfl-archive/filesystem/pkg/vfsrw"
	"github.com/ocfl-archive/filesystem/pkg/writefs"
	"github.com/ocfl-archive/filesystem/pkg/zipfs"
	"github.com/ocfl-archive/gocfl/v3/pkg/ocfl"
	"github.com/ocfl-archive/gocfl/v3/pkg/ocfl/inventory"
	"github.com/ocfl-archive/gocfl/v3/pkg/ocfl/ocflerrors"
	"github.com/ocfl-archive/gocfl/v3/pkg/ocfl/util"
	"github.com/ocfl-archive/gocfl/v3/pkg/ocfllogger"
	"github.com/ocfl-archive/indexer/v3/pkg/indexer"
)

const (
	defaultMimeType = "application/octet-stream"
	defaultPronom   = "UNKNOWN"
)

type UploaderService struct {
	StorageHandlerHandlerServiceClient handlerPb.StorageHandlerHandlerServiceClient
	ConfigObj                          config.Config
	Logger                             *ocfllogger.OCFLLogger
	Vfs                                vfsrw.VFSRW
}

func (u *UploaderService) TenantHasAccess(key string, collection string) (bool, error) {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	status, err := u.StorageHandlerHandlerServiceClient.TenantHasAccess(ctx, &pb.UploaderAccessObject{Key: key, Collection: collection})
	if err != nil {
		return false, errors.Wrapf(err, "could not get tenant access status for tenant with key: %v", key)
	}
	return status.Ok, nil
}

func (u *UploaderService) StoringFiles(order *pb.IncomingOrder, partitionId string, severalObjects string) error {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	_, err := StoringFiles(u.StorageHandlerHandlerServiceClient, ctx, order, partitionId, severalObjects)
	if err != nil {
		return errors.Wrapf(err, "cannot copy file for collection '%s'", order.CollectionAlias)
	}

	_, err = u.StorageHandlerHandlerServiceClient.AlterStatus(ctx, &pb.StatusObject{Id: order.StatusId, Status: "archived"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", order.CollectionAlias)
	}
	basepath, _ := strings.CutSuffix(order.FilePath, order.FileName)
	nName := strings.TrimSuffix(order.FileName, filepath.Ext(order.FilePath))
	filePaths := []string{order.FilePath + ".info", basepath + nName + ".json"}
	_, err = DeleteTemporaryFiles(filePaths, u.Vfs, *u.Logger)
	if err != nil {
		return errors.Wrapf(err, "cannot delete temporary files for collection '%s'", order.CollectionAlias)
	}

	return nil
}

func (u *UploaderService) CreateObjectAndFiles(tusePath string, object models.Object, collectionAlias string, basePathString string, severalObjects string) (*pb.ObjectAndFiles, error) {
	var err error
	var fileObjects []*pb.File
	head := "v1"
	versions := "{\"v1\" : {}}"
	if !object.Binary {
		fileObjects, head, versions, err = extractMetadata(tusePath, basePathString, u.Vfs, *u.Logger)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot ExtractMetadata for: %s", tusePath)
		}
	} else if severalObjects == "1" { // if object has index 1, which means it is the second object and first was json file with files, with same name but json extension
		fileObjects, err = GetFilesFromGocflObject(tusePath, basePathString, u.Vfs, *u.Logger)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot GetFilesFromGocflObject for: %s", tusePath)
		}
	}
	object.Head = head
	object.Versions = versions
	objectPb := mapper.ConvertToObjectPb(object)

	objectAndFiles := &pb.ObjectAndFiles{CollectionAlias: collectionAlias, Object: objectPb, Files: fileObjects}
	return objectAndFiles, nil
}

func extractMetadata(tusFileName string, basePath string, vfs vfsrw.VFSRW, logger ocfllogger.OCFLLogger) ([]*pb.File, string, string, error) {

	ocflPath := path.Join(basePath, tusFileName)
	destFS, err := zipfs.NewFSFile(vfs, ocflPath, logger.Logger())
	if err != nil {
		logger.Error().Err(err).Msgf("cannot open zip filesystem for '%s'", ocflPath)
		return nil, "", "", err
	}
	defer func() {
		if err := writefs.Close(destFS); err != nil {
			logger.Error().Err(err).Msgf("cannot close filesystem for '%s'", ocflPath)
		}
	}()
	var objFsys fs.FS
	_, err = util.GetStorageRootVersion(destFS)
	if err != nil && !errors.Is(err, ocflerrors.ErrVersionNone) {
		logger.Error().Err(err).Msgf("cannot get storage root version for '%s'", ocflPath)
		return nil, "", "", err
	} else if err == nil {

		fis, err := fs.ReadDir(destFS, ".")
		if err != nil {
			logger.Error().Err(err).Msgf("cannot read directory for '%s'", ocflPath)
			return nil, "", "", err
		}
		var objF string
		for _, fi := range fis {
			if fi.IsDir() && fi.Name() != "extensions" {
				objF = fi.Name()
				break
			}
		}
		if objF == "" {
			logger.Error().Msgf("cannot find OCFL object directory for '%s'", ocflPath)
			return nil, "", "", err
		}

		objFsys, err = writefs.Sub(destFS, objF)
		if err != nil {
			logger.Error().Err(err).Msgf("cannot open filesystem for '%s'", ocflPath)
			return nil, "", "", err
		}
	}

	objLoaded, err := ocfl.LoadObject(context.Background(), objFsys, nil, logger)
	if err != nil {
		logger.Error().Msgf("failed to load object '%s' at '%s': %v", tusFileName, ocflPath, err)
		return nil, "", "", err
	}
	defer objLoaded.Close()

	extractor := objLoaded.GetExtractor()
	defer func() { _ = extractor.Close() }()
	metadata, err := extractor.GetMetadata()
	if err != nil {
		logger.Error().Msgf("failed to get metadata for object '%s': %v", ocflPath, err)
		return nil, "", "", err
	}

	filesRetrieved := metadata.Files
	head := metadata.Head.String()
	versionsMap := metadata.Versions
	versionsJson, err := json.Marshal(versionsMap)
	if err != nil {
		fmt.Printf("cannot marshal versions to Json from storage root: %v\n", err)
		logger.Error().Msgf("cannot marshal versions to Json from storage root: %v\n", err)
		return nil, "", "", err
	}

	files := make([]*pb.File, 0)
	it := 0
	for index, fileRetr := range filesRetrieved {
		file := pb.File{}
		file.Checksum = index
		file.Name = fileRetr.VersionName[head]

		extensions := fileRetr.Extension["NNNN-indexer"]
		if extensions != nil {
			switch v := extensions.(type) {
			case *indexer.ResultV2:
				file.Size = int64(v.Size)
				file.Pronom = v.Pronom
				if file.Pronom == "" {
					file.Pronom = defaultPronom
				}
				file.Duration = int64(v.Duration)
				file.Width = int64(v.Width)
				file.Height = int64(v.Height)
				file.MimeType = v.Mimetype
				if file.MimeType == "" {
					file.MimeType = defaultMimeType
				}
				it++
			}
		} else {
			file.MimeType = defaultMimeType
			file.Pronom = defaultPronom
		}
		files = append(files, &file)
	}
	if it == 0 {
		return nil, "", "", errors.New("No files were extracted")
	}
	return files, head, string(versionsJson), nil
}

func GetFilesFromGocflObject(tusFileName string, basePathString string, vfs fs.FS, logger ocfllogger.OCFLLogger) ([]*pb.File, error) {
	var objectOcfl inventory.Metadata
	pathTus := path.Join(basePathString, strings.TrimSuffix(tusFileName, filepath.Ext(tusFileName))+".json")
	sourceFP, err := vfs.Open(pathTus)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := sourceFP.Close(); err != nil {
			logger.Error().Msgf("cannot close source: %v", err)
		}
	}()
	jsonObject, err := io.ReadAll(sourceFP)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonObject, &objectOcfl)
	if err != nil {
		logger.Error().Msgf("cannot Unmarshal jsonObject: %v", err)
		return nil, err
	}
	if objectOcfl.ID == "" {
		logger.Error().Msgf("Error mapping json: %s", pathTus)
		return nil, errors.New(fmt.Sprintf("Error mapping json: %s", pathTus))
	}
	filesRetrieved := objectOcfl.Files
	head := objectOcfl.Head.String()

	files := make([]*pb.File, 0)
	for _, fileRetr := range filesRetrieved {
		file := pb.File{}
		file.Name = fileRetr.VersionName[head]

		if fileRetr.Extension["NNNN-indexer"] != nil {
			extensions := fileRetr.Extension["NNNN-indexer"].(map[string]any)

			file.Pronom = extensions["pronom"].(string)
			if file.Pronom == "" {
				file.Pronom = defaultPronom
			}
			if extensions["size"] != nil {
				file.Size = int64(extensions["size"].(float64))
			}
			if extensions["duration"] != nil {
				file.Duration = int64(extensions["duration"].(float64))
			}
			if extensions["width"] != nil {
				file.Width = int64(extensions["width"].(float64))
			}
			if extensions["height"] != nil {
				file.Height = int64(extensions["height"].(float64))
			}
			file.MimeType = extensions["mimetype"].(string)
			if file.MimeType == "" {
				file.MimeType = defaultMimeType
			}
		} else {
			file.MimeType = defaultMimeType
			file.Pronom = defaultPronom
		}
		files = append(files, &file)
	}
	return files, nil
}
