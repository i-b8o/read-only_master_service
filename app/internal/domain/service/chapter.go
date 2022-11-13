package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type ChapterStorage interface {
	Create(ctx context.Context, chapter entity.Chapter) (uint64, error)
	DeleteAll(ctx context.Context, ID uint64) error
	GetAllIds(ctx context.Context, ID uint64) ([]uint64, error)
	GetRegulationIdByChapterId(ctx context.Context, ID uint64) (uint64, error)
}

type chapterService struct {
	storage ChapterStorage
}

func NewChapterService(storage ChapterStorage) *chapterService {
	return &chapterService{storage: storage}
}

func (s chapterService) Create(ctx context.Context, chapter entity.Chapter) (uint64, error) {
	return s.storage.Create(ctx, chapter)
}

func (s chapterService) DeleteAll(ctx context.Context, ID uint64) error {
	return s.storage.DeleteAll(ctx, ID)
}

func (s chapterService) GetAllIds(ctx context.Context, ID uint64) ([]uint64, error) {
	return s.storage.GetAllIds(ctx, ID)
}

func (s chapterService) GetRegulationIdByChapterId(ctx context.Context, ID uint64) (uint64, error) {
	return s.storage.GetRegulationIdByChapterId(ctx, ID)
}
