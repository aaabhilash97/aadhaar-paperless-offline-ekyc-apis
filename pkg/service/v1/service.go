package service

import (
	"errors"

	api "github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/pkg/api/v1"
	"go.uber.org/zap"
)

type AadhaarCacheStore interface {
	GetSession(string) (string, error)
	SaveSession(string) (string, error)
	SaveData(string, interface{}) error
	GetData(string, interface{}) error
	IsNotFoundError(error) bool
}

type AadhaarService struct {
	api.UnimplementedAadhaarServiceServer
	log               *zap.Logger
	aadhaarCacheStore AadhaarCacheStore
}

type NewAadhaarServiceI interface {
	GetLogger() *zap.Logger
	GetAadhaarCacheStore() AadhaarCacheStore
}

func NewAadhaarService(opt NewAadhaarServiceI) (api.AadhaarServiceServer, error) {
	srv := AadhaarService{}

	if opt == nil {
		return srv, errors.New("opt can't be nil")
	}
	if log := opt.GetLogger(); log == nil {
		return srv, errors.New("opt.GetLogger() can't be nil")
	} else {
		srv.log = log
	}

	if aadhaarCacheStore := opt.GetAadhaarCacheStore(); aadhaarCacheStore == nil {
		return srv, errors.New("opt.GetAadhaarCacheStore() can't be nil")
	} else {
		srv.aadhaarCacheStore = aadhaarCacheStore
	}

	return srv, nil
}
