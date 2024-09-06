package service

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	configuration "github.com/aws/aws-sdk-go-v2/config"
	"github.com/je4/filesystem/v2/pkg/s3fsrw"
	"github.com/je4/filesystem/v2/pkg/writefs"
	"github.com/je4/filesystem/v2/pkg/zipfs"
	ironmaiden "github.com/je4/indexer/v2/pkg/indexer"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerPb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/dlza-manager/mapper"
	"github.com/ocfl-archive/dlza-manager/models"
	"github.com/ocfl-archive/gocfl/v2/gocfl/cmd"
	"github.com/ocfl-archive/gocfl/v2/pkg/ocfl"
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
	daLogger := zLogger.NewZWrapper(*u.Logger)
	_, err := u.StorageHandlerHandlerServiceClient.AlterStatus(ctx, &pb.StatusObject{Id: order.StatusId, Status: "archiving"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", order.CollectionAlias)
	}
	_, err = CopyFiles(u.StorageHandlerHandlerServiceClient, ctx, order, u.ConfigObj, daLogger)
	if err != nil {
		return errors.Wrapf(err, "cannot copy file for collection '%s'", order.CollectionAlias)
	}

	_, err = u.StorageHandlerHandlerServiceClient.AlterStatus(ctx, &pb.StatusObject{Id: order.StatusId, Status: "archived"})
	if err != nil {
		return errors.Wrapf(err, "cannot set status to copy file for collection '%s'", order.CollectionAlias)
	}
	_, err = DeleteTemporaryFiles(order, u.ConfigObj, daLogger)
	if err != nil {
		return errors.Wrapf(err, "cannot delete temporary files for collection '%s'", order.CollectionAlias)
	}

	return nil
}

func (u *UploaderService) CreateObjectAndFiles(tusePath string, objectJson string, collectionAlias string, confObj configuration.Config) (*pb.ObjectAndFiles, error) {
	object := models.Object{}
	err := json.Unmarshal([]byte(objectJson), &object)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal object: %v", objectJson)
	}
	objectPb := mapper.ConvertToObjectPb(object)

	fileObjects, err := extractMetadata(tusePath, confObj, *u.Logger)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot ExtractMetadata for: %v", tusePath)
	}
	objectAndFiles := &pb.ObjectAndFiles{CollectionAlias: collectionAlias, Object: objectPb, Files: fileObjects}
	return objectAndFiles, nil
}

func extractMetadata(tusFileName string, confObj configuration.Config, logger zLogger.ZLogger) ([]*pb.File, error) {
	daLogger := zLogger.NewZWrapper(logger)
	fsFactory, err := writefs.NewFactory()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create filesystem factory")
	}
	// arn:cache:s3:zurich:trallala
	if err := fsFactory.Register(s3fsrw.NewCreateFSFunc(map[string]*s3fsrw.S3Access{
		"cache": &s3fsrw.S3Access{
			AccessKey: "AKIAFEDBDB2704C24D21",
			SecretKey: "0jmsjtQd0ka66thzFDJn6ESUeiLii4dIHHHgTPU6",
			URL:       "vip-ecs-ub.storage.p.unibas.ch",
			UseSSL:    true,
		},
	}, s3fsrw.ARNRegexStr, false, nil, daLogger), "^arn:", writefs.LowFS); err != nil {
		return nil, errors.Wrap(err, "cannot register zipfs")
	}
	if err := fsFactory.Register(zipfs.NewCreateFSFunc(), "([0-9a-f]{32}|\\.zip)$", writefs.HighFS); err != nil {
		return nil, errors.Wrap(err, "cannot register zipfs")
	}
	/*
		if err := fsFactory.Register(osfsrw.NewCreateFSFunc(), "", writefs.LowFS); err != nil {
			return nil, errors.Wrap(err, "cannot register zipfs")
		}
	*/

	ocflFS, err := fsFactory.Get("arn:cache:s3:::" + "ubbasel-test" + "/" + tusFileName)

	if err != nil {
		daLogger.Errorf("cannot get filesystem for file '%s': %v", tusFileName, err)
		daLogger.Debugf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		return nil, err
	}
	defer func() {
		if err := writefs.Close(ocflFS); err != nil {
			daLogger.Errorf("cannot close filesystem: %v", err)
			daLogger.Errorf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		}
	}()

	extensionFactory, err := cmd.InitExtensionFactory(map[string]string{},
		"",
		false,
		nil,
		nil,
		nil,
		nil,
		logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot instantiate extension factory")
	}

	ctx := ocfl.NewContextValidation(context.TODO())
	storageRoot, err := ocfl.LoadStorageRoot(ctx, ocflFS, extensionFactory, logger)
	if err != nil {
		daLogger.Errorf("cannot open storage root: %v", err)
		daLogger.Debugf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		return nil, err
	}
	metadata, err := storageRoot.ExtractMeta("", "")
	if err != nil {
		fmt.Printf("cannot extract metadata from storage root: %v\n", err)
		daLogger.Errorf("cannot extract metadata from storage root: %v\n", err)
		daLogger.Debugf("%v%+v", err, ocfl.GetErrorStacktrace(err))
		return nil, err
	}

	object := &ocfl.ObjectMetadata{}
	for _, mapItem := range metadata.Objects {
		object = mapItem
	}
	filesRetrieved := object.Files
	head := object.Head

	files := make([]*pb.File, 0)
	for index, fileRetr := range filesRetrieved {
		file := pb.File{}
		file.Checksum = index
		file.Name = fileRetr.VersionName[head]

		extensions := fileRetr.Extension["NNNN-indexer"]
		if extensions != nil {
			switch v := extensions.(type) {
			case *ironmaiden.ResultV2:
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
			}
		} else {
			file.MimeType = defaultMimeType
			file.Pronom = defaultPronom
		}
		files = append(files, &file)
	}
	return files, nil
}
