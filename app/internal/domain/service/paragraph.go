package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type ParagraphStorage interface {
	GetAllById(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error)
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error)
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type paragraphService struct {
	storage ParagraphStorage
}

func NewParagraphService(storage ParagraphStorage) *paragraphService {
	return &paragraphService{storage: storage}
}

func (s *paragraphService) GetAllById(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error) {
	return s.storage.GetAllById(ctx, chapterID)
}

func (s *paragraphService) CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error {
	return s.storage.CreateAll(ctx, paragraphs)
}

func (s *paragraphService) UpdateOne(ctx context.Context, content string, paragraphID uint64) error {
	return s.storage.UpdateOne(ctx, content, paragraphID)
}

func (s *paragraphService) GetWithHrefs(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error) {
	return s.storage.GetWithHrefs(ctx, chapterID)
}
func (s *paragraphService) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	return s.storage.DeleteForChapter(ctx, chapterID)
}
