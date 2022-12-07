package usecase_paragraph

import (
	"context"
	"read-only_master_service/internal/domain/entity"
	"regexp"
	"strconv"
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

type LinkService interface {
	Create(ctx context.Context, link entity.Link) error
}

// TODO delete
// type SpeechService interface {
// 	Create(ctx context.Context, speech entity.Speech) (string, error)
// }

type paragraphUsecase struct {
	paragraphService ParagraphService
	chapterService   ChapterService
	linkService      LinkService
}

func NewParagraphUsecase(paragraphService ParagraphService, chapterService ChapterService, linkService LinkService) *paragraphUsecase {
	return &paragraphUsecase{paragraphService: paragraphService, chapterService: chapterService, linkService: linkService}
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
	// the regulation id will be used for a link creation
	rId, err := u.chapterService.GetRegulationIdByChapterId(ctx, paragraphs[0].ChapterID)
	if err != nil {
		return err
	}
	// create links for paragraphs
	for _, p := range paragraphs {
		if p.ID > 0 { // sometimes any paragraph can be without an id and no one will link to it
			u.linkService.Create(ctx, entity.Link{ID: p.ID, ParagraphNum: p.Num, ChapterID: p.ChapterID, RID: rId})
			// TODO delete
			// speechTextSlice, err := speech.CreateSpeechText(p.Content, 255, 40)
			// if err != nil {
			// 	return err
			// }
			// for i, text := range speechTextSlice {
			// 	speech := entity.Speech{ParagraphID: p.ID, Content: text, OrderNum: uint64(i)}
			// 	_, err := u.speechService.Create(ctx, speech)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		}

		// when the paragraph has additional IDs inside itself we need to create for them additional links
		hasIDsInside := strings.Contains(p.Content, "<a id=")
		if hasIDsInside {
			re := regexp.MustCompile(`<a id='(.*?)'`)
			matches := re.FindAllString(p.Content, -1)
			for _, match := range matches {
				// convert the id number of the ID to the uint64
				re := regexp.MustCompile(`[\d]+`)
				subIndexStr := re.FindString(match)
				subIndexUint64, err := strconv.ParseUint(subIndexStr, 10, 64)
				if err != nil {
					return err
				}
				u.linkService.Create(ctx, entity.Link{ID: subIndexUint64, ParagraphNum: p.Num, ChapterID: p.ChapterID, RID: rId})
			}
		}
		// drop unnecessary spaces from the paragraph content
		content := strings.TrimSpace(p.Content)
		re := regexp.MustCompile(`\r?\n`)
		clearContent := re.ReplaceAllString(content, " ")

		p.Content = clearContent
	}

	err = u.paragraphService.CreateAll(ctx, paragraphs)
	if err != nil {
		return err
	}
	return err
}
