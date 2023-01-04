package usecase_paragraph

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type ParagraphService interface {
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
	UpdateOne(ctx context.Context, id, chapterID uint64, content string) error
	GetOne(ctx context.Context, paragraphId, chapterID uint64) (entity.Paragraph, error)
}

type ChapterService interface {
	GetDocIdByChapterId(ctx context.Context, ID uint64) (uint64, error)
}

type paragraphUsecase struct {
	paragraphService ParagraphService
	chapterService   ChapterService
}

func NewParagraphUsecase(paragraphService ParagraphService, chapterService ChapterService) *paragraphUsecase {
	return &paragraphUsecase{paragraphService: paragraphService, chapterService: chapterService}
}

func (u paragraphUsecase) UpdateOne(ctx context.Context, id, chapterId uint64, content string) error {
	return u.paragraphService.UpdateOne(ctx, id, chapterId, content)
}
func (u paragraphUsecase) GetOne(ctx context.Context, paragraphId, chapterID uint64) (entity.Paragraph, error) {
	return u.paragraphService.GetOne(ctx, paragraphId, chapterID)
}
func (u paragraphUsecase) CreateParagraphs(ctx context.Context, paragraphs []entity.Paragraph) error {
	if len(paragraphs) == 0 {
		return nil
	}
	return u.paragraphService.CreateAll(ctx, paragraphs)
}
