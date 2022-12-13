package usecase_chapter

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type ChapterService interface {
	Create(ctx context.Context, chapter entity.Chapter) (uint64, error)
}

type PseudoChapter interface {
	CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error
	DeleteRelationship(ctx context.Context, chapterID uint64) error
}

type chapterUsecase struct {
	chapterService ChapterService

	pseudoChapter PseudoChapter
	logging       logging.Logger
}

func NewChapterUsecase(chapterService ChapterService, pseudoChapter PseudoChapter, logging logging.Logger) *chapterUsecase {
	return &chapterUsecase{chapterService: chapterService, pseudoChapter: pseudoChapter, logging: logging}
}

func (u chapterUsecase) CreateChapter(ctx context.Context, chapter entity.Chapter) (uint64, error) {
	// create a chapter
	ID, err := u.chapterService.Create(ctx, chapter)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// create an id-pseudoId relationship
	err = u.pseudoChapter.CreateRelationship(ctx, entity.PseudoChapter{ID: ID, PseudoId: chapter.Pseudo})
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}
	return ID, nil
}
