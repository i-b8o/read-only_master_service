package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type DocStorage interface {
	GetAll(ctx context.Context) ([]entity.Doc, error)
	Create(ctx context.Context, doc entity.Doc) (uint64, error)
	Delete(ctx context.Context, docID uint64) error
}

type docService struct {
	storage DocStorage
	logging logging.Logger
}

func NewDocService(storage DocStorage, logging logging.Logger) *docService {
	return &docService{storage: storage, logging: logging}
}

func (s *docService) Create(ctx context.Context, doc entity.Doc) (uint64, error) {
	id, err := s.storage.Create(ctx, doc)
	if err != nil {
		s.logging.Errorf("%v %v", doc, err)
		return 0, err
	}
	return id, nil
}

func (s *docService) Delete(ctx context.Context, docId uint64) error {
	err := s.storage.Delete(ctx, docId)
	if err != nil {
		s.logging.Errorf("%d %v", docId, err)
		return err
	}
	return nil
}

func (s *docService) GetAll(ctx context.Context) ([]entity.Doc, error) {
	docs, err := s.storage.GetAll(ctx)
	if err != nil {
		s.logging.Error(err)
		return nil, err
	}
	return docs, nil
}
