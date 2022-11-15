package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type PseudoChapterStorage interface {
	CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error
	DeleteRelationship(ctx context.Context, chapterID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type pseudoChapterService struct {
	storage PseudoChapterStorage
}

func NewPseudoChapterService(storage PseudoChapterStorage) *pseudoChapterService {
	return &pseudoChapterService{storage: storage}
}

func (pcs pseudoChapterService) CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error {
	return pcs.storage.CreateRelationship(ctx, pseudoChapter)
}

func (pcs pseudoChapterService) DeleteRelationship(ctx context.Context, chapterID uint64) error {
	return pcs.storage.DeleteRelationship(ctx, chapterID)
}

func (pcs pseudoChapterService) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	return pcs.storage.GetIDByPseudo(ctx, pseudoId)
}
