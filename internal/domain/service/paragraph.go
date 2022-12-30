package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type ParagraphStorage interface {
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
	GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error)
	UpdateOne(ctx context.Context, id uint64, content string) error
	GetOne(ctx context.Context, paragraphId, chapterID uint64) (entity.Paragraph, error)
}

type paragraphService struct {
	storage ParagraphStorage
	logging logging.Logger
}

func NewParagraphService(storage ParagraphStorage, logging logging.Logger) *paragraphService {
	return &paragraphService{storage: storage, logging: logging}
}

func (s *paragraphService) CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error {
	err := s.storage.CreateAll(ctx, paragraphs)
	if err != nil {
		s.logging.Errorf("%v %v", paragraphs, err)
		return err
	}
	return nil
}

func (s *paragraphService) GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error) {
	paragraphs, err := s.storage.GetParagraphsWithHrefs(ctx, chapterId)
	if err != nil {
		s.logging.Errorf("%d %v", chapterId, err)
		return nil, err
	}
	return paragraphs, nil
}

func (s *paragraphService) UpdateOne(ctx context.Context, id uint64, content string) error {
	err := s.storage.UpdateOne(ctx, id, content)
	if err != nil {
		s.logging.Errorf("%d %s %v", id, content, err)
		return err
	}
	return nil
}

func (s *paragraphService) GetOne(ctx context.Context, paragraphId, chapterID uint64) (*entity.Paragraph, error) {
	paragraph, err := s.storage.GetOne(ctx, paragraphId, chapterID)
	if err != nil {
		s.logging.Errorf("%d %d %v", paragraphId, chapterID, err)
		return nil, err
	}
	return &paragraph, nil
}
