package usecase_regulation

import (
	"context"
	"fmt"
	"regexp"
	"regulations_supreme_service/internal/domain/entity"
	"strings"

	"github.com/i-b8o/logging"
)

type RegulationService interface {
	Create(ctx context.Context, regulation entity.Regulation) (uint64, error)
	Delete(ctx context.Context, regulationId uint64) error
}

type ChapterService interface {
	DeleteAll(ctx context.Context, ID uint64) error
	GetAllIds(ctx context.Context, ID uint64) ([]uint64, error)
}
type ParagraphService interface {
	DeleteForRegulation(ctx context.Context, chaptersIDs []uint64) error
	GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error)
	UpdateOne(ctx context.Context, id uint64, content string) error
}

type AbsentService interface {
	Done(ctx context.Context, pseudo string) error
	Create(ctx context.Context, absent entity.Absent) error
}

type PseudoRegulationService interface {
	CreateRelationship(ctx context.Context, pseudoRegulation entity.PseudoRegulation) error
	DeleteRelationship(ctx context.Context, regulationID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type PseudoChapterService interface {
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}
type regulationUsecase struct {
	regulationService       RegulationService
	chapterService          ChapterService
	paragraphService        ParagraphService
	absentService           AbsentService
	pseudoRegulationService PseudoRegulationService
	pseudoChapterService    PseudoChapterService
	logging                 logging.Logger
}

func NewRegulationUsecase(regulationService RegulationService, chapterService ChapterService, paragraphService ParagraphService, absentService AbsentService, pseudoRegulationService PseudoRegulationService, pseudoChapterService PseudoChapterService, logging logging.Logger) *regulationUsecase {
	return &regulationUsecase{regulationService: regulationService, chapterService: chapterService, paragraphService: paragraphService, absentService: absentService, pseudoRegulationService: pseudoRegulationService, pseudoChapterService: pseudoChapterService, logging: logging}
}

func (u regulationUsecase) CreateRegulation(ctx context.Context, regulation entity.Regulation) (uint64, error) {
	// create a regulation
	ID, err := u.regulationService.Create(ctx, regulation)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// create an id-pseudoId relationship
	err = u.pseudoRegulationService.CreateRelationship(ctx, entity.PseudoRegulation{ID: ID, PseudoId: regulation.Pseudo})
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// mark the regulation as done
	err = u.absentService.Done(ctx, regulation.Pseudo)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	return ID, nil
}

func (u regulationUsecase) DeleteRegulation(ctx context.Context, ID uint64) error {
	// delete a regulation
	err := u.regulationService.Delete(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}

	// delete all paragraphs for the regulation
	IDs, err := u.chapterService.GetAllIds(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	err = u.paragraphService.DeleteForRegulation(ctx, IDs)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	// delete all chapters for the regulation
	err = u.chapterService.DeleteAll(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}

	// delete the id-pseudoId relationship
	err = u.pseudoRegulationService.DeleteRelationship(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	return nil
}

func (u regulationUsecase) GenerateLinks(ctx context.Context, regulationID uint64) error {

	// for every chapter in the regulation
	chIDs, err := u.chapterService.GetAllIds(ctx, regulationID)
	if err != nil {
		return err
	}

	for _, chId := range chIDs {
		// get only paragraphs with links inside
		paragraphs, err := u.paragraphService.GetParagraphsWithHrefs(ctx, chId)
		if err != nil {
			u.logging.Error(err)
			return err
		}

		for _, paragraph := range paragraphs {
			content := paragraph.Content
			re := regexp.MustCompile("<a href='.+'>")
			links := re.FindAllString(content, -1)

			for _, aLink := range links {
				hrefRaw := strings.Split(aLink, "<a href='")[1]
				href := strings.Split(hrefRaw, "'>")[0]
				// link for a entire document
				matched, err := regexp.MatchString(`^\/document\/cons_doc_LAW_\d+\/$`, href)
				if err != nil {
					u.logging.Error(err)
				}
				if matched {
					// fetch only doc ID (pseudo)
					p := strings.Split(href, "/document/cons_doc_LAW_")[1]
					ID := strings.Split(p, "/")[0]
					// get relative ID
					rID, err := u.pseudoRegulationService.GetIDByPseudo(ctx, ID)
					if err != nil {
						fmt.Printf("href = %s,ID: %s, error: %s", href, ID, err.Error())
					}
					// if id was not found - absent
					if rID == 0 {
						absent := entity.Absent{Pseudo: ID, ParagraphID: paragraph.ID}
						err := u.absentService.Create(ctx, absent)
						if err != nil {
							fmt.Println(err.Error())
						}
						continue
					}
				}

				// get ids for a regulation, a chapter and a paragraph
				IDs := strings.Split(href, "/")

				// something wrong with the link - absent
				if (len(IDs) < 3) || (len(IDs[0]) == 0) || (len(IDs[1]) == 0) || (len(IDs[2]) == 0) {
					absent := entity.Absent{Pseudo: href, ParagraphID: paragraph.ID}
					u.absentService.Create(ctx, absent)
					continue
				}

				// get relative ID
				regID, err := u.pseudoRegulationService.GetIDByPseudo(ctx, IDs[0])
				if err != nil {
					fmt.Printf("href = %s,ID: %s, error: %s", href, IDs[0], err.Error())
				}
				// if id was not found - absent
				if regID == 0 {
					absent := entity.Absent{Pseudo: IDs[0], ParagraphID: paragraph.ID}
					u.absentService.Create(ctx, absent)
					continue
				}
				chID, err := u.pseudoChapterService.GetIDByPseudo(ctx, IDs[1])

				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				pID := IDs[2]

				post := fmt.Sprintf("%d/%d/%s>", regID, chID, pID)
				content = strings.Replace(content, aLink, "<a href='"+post, 1)
			}
			err := u.paragraphService.UpdateOne(ctx, paragraph.ID, content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
