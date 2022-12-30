package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type PseudoChapterStorage interface {
	CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error
	DeleteRelationship(ctx context.Context, chapterID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type pseudoChapterService struct {
	storage PseudoChapterStorage
	logging logging.Logger
}

func NewPseudoChapterService(storage PseudoChapterStorage, logging logging.Logger) *pseudoChapterService {
	return &pseudoChapterService{storage: storage, logging: logging}
}

func (pcs pseudoChapterService) CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error {
	err := pcs.storage.CreateRelationship(ctx, pseudoChapter)
	if err != nil {
		pcs.logging.Errorf("%v %v", pseudoChapter, err)
		return err
	}
	return nil
}

func (pcs pseudoChapterService) DeleteRelationship(ctx context.Context, chapterID uint64) error {
	err := pcs.storage.DeleteRelationship(ctx, chapterID)
	if err != nil {
		pcs.logging.Errorf("%d %v", chapterID, err)
		return err
	}
	return nil
}

func (pcs pseudoChapterService) GetIDByPseudo(ctx context.Context, pseudoId string) (*uint64, error) {
	id, err := pcs.storage.GetIDByPseudo(ctx, pseudoId)
	if err != nil {
		pcs.logging.Errorf("%d %v", id, err)
		return nil, err
	}
	return &id, nil
}
