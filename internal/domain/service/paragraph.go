package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type ParagraphStorage interface {
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
	GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error)
	UpdateOne(ctx context.Context, id uint64, content string) error
}

type paragraphService struct {
	storage ParagraphStorage
}

func NewParagraphService(storage ParagraphStorage) *paragraphService {
	return &paragraphService{storage: storage}
}

func (s *paragraphService) DeleteForRegulation(ctx context.Context, chaptersIDs []uint64) error {
	for _, ID := range chaptersIDs {
		err := s.storage.DeleteForChapter(ctx, ID)
		if err != nil {
			return err
		}
	}
	return nil
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
