package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type DocStorage interface {
	GetAll(ctx context.Context) ([]entity.Doc, error)
	Create(ctx context.Context, doc entity.Doc) (uint64, error)
	Delete(ctx context.Context, docID uint64) error
}

type docService struct {
	storage DocStorage
}

func NewDocService(storage DocStorage) *docService {
	return &docService{storage: storage}
}

func (s *docService) Create(ctx context.Context, doc entity.Doc) (uint64, error) {
	return s.storage.Create(ctx, doc)
}

func (s *docService) Delete(ctx context.Context, docId uint64) error {
	return s.storage.Delete(ctx, docId)
}

func (s *docService) GetAll(ctx context.Context) ([]entity.Doc, error) {
	return s.storage.GetAll(ctx)
}
