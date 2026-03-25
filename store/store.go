package store

import (
	"context"
	"strings"

	"emperror.dev/errors"
	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
)

type AutoloadStore func(tenantAlias string, partitionId string) (s3store.S3Store, bool)

type RoutingStore struct {
	autoload AutoloadStore
}

func storeFromID(id string) (string, string) {
	ret := strings.Split(id, "-")
	partitionId := strings.Join(ret[1:len(ret)-1], "-")
	tenantAlias := ret[0]
	if strings.Contains(tenantAlias, "/") {
		tenantAlias = strings.Split(tenantAlias, "/")[1]
	}
	if len(ret) == 7 {
		return tenantAlias, partitionId
	}
	return "", ""
}

func (s RoutingStore) storeFromID(id string) (s3store.S3Store, error) {
	tenantAlias, partitionId := storeFromID(id)
	if tenantAlias == "" || partitionId == "" {
		return s3store.S3Store{}, errors.Wrapf(handler.ErrNotFound, "invalid upload ID: %s", id)
	}
	store, ok := s.autoload(tenantAlias, partitionId)
	if !ok {
		return s3store.S3Store{}, errors.Wrapf(handler.ErrNotFound, "invalid upload ID: %s", id)
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

func NewRoutingStore(autoload AutoloadStore) RoutingStore {
	return RoutingStore{
		autoload: autoload,
	}
}

func (s RoutingStore) AsServableUpload(upload handler.Upload) handler.ServableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		return nil
	}
	return store.AsServableUpload(upload)
}

func (s RoutingStore) AsLengthDeclarableUpload(upload handler.Upload) handler.LengthDeclarableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		return nil
	}
	return store.AsLengthDeclarableUpload(upload)
}

func (s RoutingStore) AsConcatableUpload(upload handler.Upload) handler.ConcatableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		return nil
	}
	return store.AsConcatableUpload(upload)
}

func (s RoutingStore) AsTerminatableUpload(upload handler.Upload) handler.TerminatableUpload {
	info, err := upload.GetInfo(context.Background())
	if err != nil {
		return nil
	}
	store, err := s.storeFromID(info.ID)
	if err != nil {
		return nil
	}
	return store.AsTerminatableUpload(upload)
}

func (s RoutingStore) NewUpload(ctx context.Context, info handler.FileInfo) (upload handler.Upload, err error) {
	store, err := s.storeFromID(info.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "storage is not found for id %s", info.ID)
	}
	ret, err := store.NewUpload(ctx, info)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

func (s RoutingStore) GetUpload(ctx context.Context, id string) (upload handler.Upload, err error) {
	store, err := s.storeFromID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "get upload %s", id)
	}
	ret, err := store.GetUpload(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret, nil
}
