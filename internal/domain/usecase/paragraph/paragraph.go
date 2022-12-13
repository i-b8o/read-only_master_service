package usecase_paragraph

import (
	"context"
	"read-only_master_service/internal/domain/entity"
	"regexp"
	"strings"
)

type ParagraphService interface {
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
	UpdateOne(ctx context.Context, id uint64, content string) error
	GetOne(ctx context.Context, paragraphId uint64) (entity.Paragraph, error)
}

type ChapterService interface {
	GetRegulationIdByChapterId(ctx context.Context, ID uint64) (uint64, error)
}

type paragraphUsecase struct {
	paragraphService ParagraphService
	chapterService   ChapterService
}

func NewParagraphUsecase(paragraphService ParagraphService, chapterService ChapterService) *paragraphUsecase {
	return &paragraphUsecase{paragraphService: paragraphService, chapterService: chapterService}
}

func (u paragraphUsecase) UpdateOne(ctx context.Context, id uint64, content string) error {
	return u.paragraphService.UpdateOne(ctx, id, content)
}
func (u paragraphUsecase) GetOne(ctx context.Context, paragraphId uint64) (entity.Paragraph, error) {
	return u.paragraphService.GetOne(ctx, paragraphId)
}
func (u paragraphUsecase) CreateParagraphs(ctx context.Context, paragraphs []entity.Paragraph) error {
	if len(paragraphs) == 0 {
		return nil
	}

	for _, p := range paragraphs {
		// drop unnecessary spaces from the paragraph content
		content := strings.TrimSpace(p.Content)
		re := regexp.MustCompile(`\r?\n`)
		clearContent := re.ReplaceAllString(content, " ")

		p.Content = clearContent
	}
	return u.paragraphService.CreateAll(ctx, paragraphs)
}
