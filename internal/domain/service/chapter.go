package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type ChapterStorage interface {
	Create(ctx context.Context, chapter entity.Chapter) (uint64, error)
	GetAllIds(ctx context.Context, ID uint64) ([]uint64, error)
	GetDocIdByChapterId(ctx context.Context, ID uint64) (uint64, error)
}

type chapterService struct {
	storage ChapterStorage
	logging logging.Logger
}

func NewChapterService(storage ChapterStorage, logging logging.Logger) *chapterService {
	return &chapterService{storage: storage, logging: logging}
}

func (s chapterService) Create(ctx context.Context, chapter entity.Chapter) (uint64, error) {
	id, err := s.storage.Create(ctx, chapter)
	if err != nil {
		s.logging.Errorf("%v %v", chapter, err)
		return 0, err
	}
	return id, nil
}
func (s chapterService) GetAllIds(ctx context.Context, ID uint64) ([]uint64, error) {
	IDs, err := s.storage.GetAllIds(ctx, ID)
	if err != nil {
		s.logging.Errorf("%d %v", ID, err)
		return nil, err
	}
	return IDs, nil
}

func (s chapterService) GetDocIdByChapterId(ctx context.Context, ID uint64) (uint64, error) {
	id, err := s.storage.GetDocIdByChapterId(ctx, ID)
	if err != nil {
		s.logging.Errorf("%d %v", id, err)
		return 0, err
	}
	return id, nil
}
