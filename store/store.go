package store

import (
	"context"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
)

type AutoloadStore func(tenantAlias string, partitionId string) (s3store.S3Store, error)

type RoutingStore struct {
	autoload AutoloadStore
	logger   zLogger.ZLogger
}

var storeFromIDUUIDRe = regexp.MustCompile(
	`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`,
)

func storeFromID(id string) (string, string) {
	if id == "" {
		return "", ""
	}

	loc := storeFromIDUUIDRe.FindStringIndex(id)
	if loc == nil {
		return "", ""
	}

	partitionId := id[loc[0]:loc[1]]

	// Tenant alias is everything right before the UUID (can contain '-') and may be prefixed by a path.
	tenantPart := id[:loc[0]]
	if i := strings.LastIndex(tenantPart, "/"); i >= 0 {
		tenantPart = tenantPart[i+1:]
	}
	tenantAlias := strings.TrimSuffix(tenantPart, "-")
	if tenantAlias == "" {
		return "", ""
	}

	return tenantAlias, partitionId

}

func (s RoutingStore) storeFromID(id string) (s3store.S3Store, error) {
	tenantAlias, partitionId := storeFromID(id)
	if tenantAlias == "" || partitionId == "" {
		s.logger.Error().Msgf("invalid upload ID: %s", id)
		return s3store.S3Store{}, errors.Wrapf(handler.ErrNotFound, "invalid upload ID: %s", id)
	}
	store, err := s.autoload(tenantAlias, partitionId)
	if err != nil {
		s.logger.Error().Msgf("store for tenant %s and partition %s not found", tenantAlias, partitionId)
		return s3store.S3Store{}, errors.Errorf("store for tenant %s and partition %s not found", tenantAlias, partitionId)
	}
	return store, nil
}

func (s RoutingStore) UseIn(composer *handler.StoreComposer) {
	composer.UseCore(s)
	composer.UseTerminater(s)
	composer.UseConcater(s)
	composer.UseLengthDeferrer(s)
	composer.UseContentServer(s)
}

func NewRoutingStore(autoload AutoloadStore, logger zLogger.ZLogger) RoutingStore {
	return RoutingStore{
		autoload: autoload,
		logger:   logger,
	}
}

func (s RoutingStore) AsServableUpload(upload handler.Upload) handler.ServableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		s.logger.Error().Msgf("cannot get info for upload for AsServableUpload. err: %v, upload: %v", err, upload)
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		s.logger.Error().Msgf("cannot get store for id %s", info.ID)
		return nil
	}
	return store.AsServableUpload(upload)
}

func (s RoutingStore) AsLengthDeclarableUpload(upload handler.Upload) handler.LengthDeclarableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		s.logger.Error().Msgf("cannot get info for upload for AsLengthDeclarableUpload. err: %v, upload: %v", err, upload)
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		s.logger.Error().Msgf("cannot get store for id %s", info.ID)
		return nil
	}
	return store.AsLengthDeclarableUpload(upload)
}

func (s RoutingStore) AsConcatableUpload(upload handler.Upload) handler.ConcatableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		s.logger.Error().Msgf("cannot get info for upload for AsConcatableUpload. err: %v, upload: %v", err, upload)
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		s.logger.Error().Msgf("cannot get store for id %s", info.ID)
		return nil
	}
	return store.AsConcatableUpload(upload)
}

func (s RoutingStore) AsTerminatableUpload(upload handler.Upload) handler.TerminatableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		s.logger.Error().Msgf("cannot get info for AsTerminatableUpload. err: %v, upload: %v", err, upload)
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		s.logger.Error().Msgf("cannot get store for id %s", info.ID)
		return nil
	}
	return store.AsTerminatableUpload(upload)
}

func (s RoutingStore) NewUpload(ctx context.Context, info handler.FileInfo) (upload handler.Upload, err error) {
	store, err := s.storeFromID(info.ID)
	if err != nil {
		s.logger.Error().Msgf("cannot get store for id %s", info.ID)
		return nil, errors.Wrapf(err, "storage is not found for id %s", info.ID)
	}
	ret, err := store.NewUpload(ctx, info)
	if err != nil {
		s.logger.Error().Msgf("cannot create new upload. err: %v, info: %v", err, info)
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

func (s RoutingStore) GetUpload(ctx context.Context, id string) (upload handler.Upload, err error) {
	store, err := s.storeFromID(id)
	if err != nil {
		s.logger.Error().Msgf("cannot get store for id %s", id)
		return nil, errors.Wrapf(err, "get upload %s", id)
	}
	ret, err := store.GetUpload(ctx, id)
	if err != nil {
		s.logger.Error().Msgf("cannot get upload %s. err: %v", id, err)
		return nil, errors.WithStack(err)
	}
	return ret, nil
}
