package usecase_chapter

import (
	"context"
	"fmt"
	"regulations_supreme_service/internal/domain/entity"
	"strconv"
)

type ChapterService interface {
	Create(ctx context.Context, chapter entity.Chapter) (string, error)
	// GetOneById(ctx context.Context, chapterID uint64) (entity.Chapter, error)
	// GetAllById(ctx context.Context, regulationID uint64) ([]entity.Chapter, error)
}

// type ParagraphService interface {
// GetAllById(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error)
// }

// type RegulationService interface {
// GetOne(ctx context.Context, regulationID uint64) (entity.Regulation, error)
// }

type LinkService interface {
	CreateForChapter(ctx context.Context, link entity.Link) error
}

type chapterUsecase struct {
	chapterService ChapterService
	// paragraphService  ParagraphService
	linkService LinkService
	// regulationService RegulationService
}

func NewChapterUsecase(chapterService ChapterService, linkService LinkService) *chapterUsecase {
	return &chapterUsecase{chapterService: chapterService, linkService: linkService}
}

func (u chapterUsecase) CreateChapter(ctx context.Context, chapter entity.Chapter) string {
	id, err := u.chapterService.Create(ctx, chapter)
	if err != nil {
		fmt.Printf("Create %s", err.Error())
		return ""
	}

	if chapter.ID > 0 {
		ch_id, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			fmt.Printf("ParseUint %s", err.Error())
			return ""
		}

		err = u.linkService.CreateForChapter(ctx, entity.Link{ID: chapter.ID, ParagraphNum: 0, ChapterID: ch_id, RID: chapter.RegulationID})
		if err != nil {
			fmt.Printf("CreateForChapter %s", err.Error())
			return ""
		}
	}
	return id
}

// func (u chapterUsecase) GetChapter(ctx context.Context, chapterID string) (entity.Regulation, entity.Chapter) {
// 	uint64ID, err := strconv.ParseUint(chapterID, 10, 64)
// 	if err != nil {
// 		return entity.Regulation{}, entity.Chapter{}
// 	}

// 	chapter, err := u.chapterService.GetOneById(ctx, uint64ID)
// 	if err != nil {
// 		return entity.Regulation{}, entity.Chapter{}
// 	}

// 	chapter.Paragraphs, err = u.paragraphService.GetAllById(ctx, uint64ID)
// 	if err != nil {
// 		return entity.Regulation{}, entity.Chapter{}
// 	}
// 	regulation, err := u.regulationService.GetOne(ctx, chapter.RegulationID)
// 	if err != nil {
// 		return entity.Regulation{}, entity.Chapter{}
// 	}
// 	chapters, err := u.chapterService.GetAllById(ctx, chapter.RegulationID)
// 	if err != nil {
// 		return entity.Regulation{}, entity.Chapter{}
// 	}
// 	regulation.Chapters = chapters
// 	return regulation, chapter
// }
