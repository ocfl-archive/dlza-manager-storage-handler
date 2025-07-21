package service

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/je4/filesystem/v3/pkg/s3fsrw"
	"github.com/je4/filesystem/v3/pkg/writefs"
	"github.com/je4/filesystem/v3/pkg/zipfs"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/dlza-manager/mapper"
	"github.com/ocfl-archive/dlza-manager/models"
	archiveerror "github.com/ocfl-archive/error/pkg/error"
	"github.com/ocfl-archive/gocfl/v2/gocfl/cmd"
	"github.com/ocfl-archive/gocfl/v2/pkg/ocfl"
	"github.com/ocfl-archive/indexer/v3/pkg/indexer"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMimeType = "application/octet-stream"
	defaultPronom   = "UNKNOWN"
)

type UploaderService struct {
	StorageHandlerHandlerServiceClient handlerPb.StorageHandlerHandlerServiceClient
	ConfigObj                          config.Config
	Logger                             *zLogger.ZLogger
	Vfs                                fs.FS
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

func (u *UploaderService) CopyFiles(order *pb.IncomingOrder) error {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	_, err := u.StorageHandlerHandlerServiceClient.AlterStatus(ctx, &pb.StatusObject{Id: order.StatusId, Status: "archiving"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", order.CollectionAlias)
	}
	_, err = CopyFiles(u.StorageHandlerHandlerServiceClient, ctx, order, u.Vfs, *u.Logger)
	if err != nil {
		return errors.Wrapf(err, "cannot copy file for collection '%s'", order.CollectionAlias)
	}

	_, err = u.StorageHandlerHandlerServiceClient.AlterStatus(ctx, &pb.StatusObject{Id: order.StatusId, Status: "archived"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", order.CollectionAlias)
	}
	_, err = DeleteTemporaryFiles(order.FilePath, u.Vfs, *u.Logger)
	if err != nil {
		return errors.Wrapf(err, "cannot delete temporary files for collection '%s'", order.CollectionAlias)
	}

	return nil
}

func (u *UploaderService) CreateObjectAndFiles(tusePath string, objectJson string, collectionAlias string, basePathString string, severalObjects string, confObj config.Config, errorFactory *archiveerror.Factory) (*pb.ObjectAndFiles, error) {
	object := models.Object{}
	err := json.Unmarshal([]byte(objectJson), &object)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal object: %s", objectJson)
	}
	var fileObjects []*pb.File
	head := "v1"
	versions := "{\"v1\" : {}}"
	if !object.Binary {
		fileObjects, head, versions, err = extractMetadata(tusePath, confObj, *u.Logger, errorFactory)
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

func extractMetadata(tusFileName string, conf config.Config, logger zLogger.ZLogger, errorFactory *archiveerror.Factory) ([]*pb.File, string, string, error) {
	fsFactory, err := writefs.NewFactory()
	if err != nil {
		return nil, "", "", errors.Wrap(err, "cannot create filesystem factory")
	}
	// arn:cache:s3:zurich:trallala
	if err := fsFactory.Register(s3fsrw.NewCreateFSFunc(map[string]*s3fsrw.S3Access{
		"cache": &s3fsrw.S3Access{
			AccessKey: conf.S3TempStorage.Key,
			SecretKey: conf.S3TempStorage.Secret,
			URL:       conf.S3TempStorage.Url,
			UseSSL:    true,
		},
	}, s3fsrw.ARNRegexStr, false, nil, "", "", logger), "^arn:", writefs.LowFS); err != nil {
		return nil, "", "", errors.Wrap(err, "cannot register zipfs")
	}
	if err := fsFactory.Register(zipfs.NewCreateFSFunc(logger), "([0-9a-f]{32}|\\.zip)$", writefs.HighFS); err != nil {
		return nil, "", "", errors.Wrap(err, "cannot register zipfs")
	}
	/*
		if err := fsFactory.Register(osfsrw.NewCreateFSFunc(), "", writefs.LowFS); err != nil {
			return nil, errors.Wrap(err, "cannot register zipfs")
		}
	*/

	ocflFS, err := fsFactory.Get("arn:cache:s3:::"+conf.S3TempStorage.Bucket+"/"+tusFileName, true)
	if err != nil {
		logger.Error().Msgf("cannot get filesystem for file '%s': %v", tusFileName, err)
		logger.Debug().Msgf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		return nil, "", "", err
	}
	defer func() {
		if err := writefs.Close(ocflFS); err != nil {
			logger.Error().Msgf("cannot close filesystem: %v", err)
			logger.Error().Msgf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		}
	}()

	extensionFactory, err := cmd.InitExtensionFactory(map[string]string{},
		"",
		false,
		nil,
		nil,
		nil,
		nil,
		logger,
		"")
	if err != nil {
		return nil, "", "", errors.Wrap(err, "cannot instantiate extension factory")
	}

	ctx := ocfl.NewContextValidation(context.TODO())
	storageRoot, err := ocfl.LoadStorageRoot(ctx, ocflFS, extensionFactory, logger, errorFactory, "")
	if err != nil {
		logger.Error().Msgf("cannot open storage root: %v", err)
		logger.Debug().Msgf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		return nil, "", "", err
	}
	metadata, err := storageRoot.ExtractMeta("", "")
	if err != nil {
		fmt.Printf("cannot extract metadata from storage root: %v\n", err)
		logger.Error().Msgf("cannot extract metadata from storage root: %v\n", err)
		logger.Debug().Msgf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		return nil, "", "", err
	}

	object := &ocfl.ObjectMetadata{}
	for _, mapItem := range metadata.Objects {
		object = mapItem
	}
	filesRetrieved := object.Files
	head := object.Head
	versionsMap := object.Versions
	versionsJson, err := json.Marshal(versionsMap)
	if err != nil {
		fmt.Printf("cannot marshal versions to Json from storage root: %v\n", err)
		logger.Error().Msgf("cannot marshal versions to Json from storage root: %v\n", err)
		logger.Debug().Msgf("%v%+v", err, ocfl.GetErrorStacktrace(err))
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

func GetFilesFromGocflObject(tusFileName string, basePathString string, vfs fs.FS, logger zLogger.ZLogger) ([]*pb.File, error) {
	objectOcfl := ocfl.StorageRootMetadata{}
	object := &ocfl.ObjectMetadata{}
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
		logger.Error().Msgf(err.Error())
		return nil, err
	}
	for _, mapItem := range objectOcfl.Objects {
		object = mapItem
	}
	filesRetrieved := object.Files
	head := object.Head

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
