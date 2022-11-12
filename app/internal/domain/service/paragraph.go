package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type ParagraphStorage interface {
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
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
