package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type ParagraphStorage interface {
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
	GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error)
	UpdateOne(ctx context.Context, id uint64, content string) error
	GetOne(ctx context.Context, paragraphId, chapterID uint64) (entity.Paragraph, error)
}

type paragraphService struct {
	storage ParagraphStorage
}

func NewParagraphService(storage ParagraphStorage) *paragraphService {
	return &paragraphService{storage: storage}
}

func (s *paragraphService) CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error {
	return s.storage.CreateAll(ctx, paragraphs)
}

func (s *paragraphService) GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error) {
	return s.storage.GetParagraphsWithHrefs(ctx, chapterId)
}

func (s *paragraphService) UpdateOne(ctx context.Context, id uint64, content string) error {
	return s.storage.UpdateOne(ctx, id, content)
}

func (s *paragraphService) GetOne(ctx context.Context, paragraphId, chapterID uint64) (entity.Paragraph, error) {
	return s.storage.GetOne(ctx, paragraphId, chapterID)
}
