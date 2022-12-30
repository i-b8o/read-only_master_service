package usecase_chapter

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
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
	pseudoChapter  PseudoChapter
}

func NewChapterUsecase(chapterService ChapterService, pseudoChapter PseudoChapter) *chapterUsecase {
	return &chapterUsecase{chapterService: chapterService, pseudoChapter: pseudoChapter}
}

func (u chapterUsecase) CreateChapter(ctx context.Context, chapter entity.Chapter) (uint64, error) {
	// create a chapter
	ID, err := u.chapterService.Create(ctx, chapter)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// create an id-pseudoId relationship
	fmt.Printf("Chapter pseudoID: %s", chapter.Pseudo)
	err = u.pseudoChapter.CreateRelationship(ctx, entity.PseudoChapter{ID: ID, PseudoId: chapter.Pseudo})
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}
	return ID, nil
}
