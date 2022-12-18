package usecase_regulation

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	"regexp"
	"strings"

	"github.com/i-b8o/logging"
)

type RegulationService interface {
	GetAll(ctx context.Context) ([]entity.Regulation, error)
	Create(ctx context.Context, regulation entity.Regulation) (uint64, error)
	Delete(ctx context.Context, regulationId uint64) error
}

type ChapterService interface {
	GetAllIds(ctx context.Context, ID uint64) ([]uint64, error)
}
type ParagraphService interface {
	GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error)
	UpdateOne(ctx context.Context, id uint64, content string) error
}

type AbsentService interface {
	Done(ctx context.Context, pseudo string) error
	Create(ctx context.Context, absent entity.Absent) error
	GetAll(ctx context.Context) ([]*entity.Absent, error)
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

func (u regulationUsecase) GetAll(ctx context.Context) ([]entity.Regulation, error) {
	return u.regulationService.GetAll(ctx)
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

func (u regulationUsecase) GetAbsents(ctx context.Context) ([]*entity.Absent, error) {
	absents, err := u.absentService.GetAll(ctx)
	if err != nil {
		u.logging.Error(err)
		return nil, err
	}

	return absents, nil
}

func (u regulationUsecase) DeleteRegulation(ctx context.Context, ID uint64) error {
	// delete a regulation
	err := u.regulationService.Delete(ctx, ID)
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
	// get IDs for every chapter in the regulation
	chIDs, err := u.chapterService.GetAllIds(ctx, regulationID)
	if err != nil {
		u.logging.Error(err)
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

			// get links
			content := paragraph.Content
			re := regexp.MustCompile("<a href='.+?'>")
			links := re.FindAllString(content, -1)

			for _, aLink := range links {

				hrefRaw := strings.Split(aLink, "<a href='")[1]
				href := strings.Split(hrefRaw, "'>")[0]
				// get if exist IDs for reggulation, chapter and paragrap
				rID, chID, pID := getIDs(href)

				// something wrong with the link - absent
				if rID == "" {

					absent := entity.Absent{Pseudo: href, ParagraphID: paragraph.ID}
					u.absentService.Create(ctx, absent)
					u.logging.Error(err)
					continue
				}
				// link for an entire document.
				if chID == "" {

					// get relative regulation ID
					regulationID, err := u.pseudoRegulationService.GetIDByPseudo(ctx, rID)
					if err != nil {
						u.logging.Error(err)
					}
					if regulationID == 0 {
						absent := entity.Absent{Pseudo: rID, ParagraphID: paragraph.ID}
						err := u.absentService.Create(ctx, absent)
						if err != nil {
							u.logging.Error(err)
						}
						continue
					}
					post := fmt.Sprintf("%d'>", regulationID)
					content = strings.Replace(content, aLink, "<a href='"+post, 1)
				}

				// link for a paragraph
				// get relative regulation ID
				regulationID, err := u.pseudoRegulationService.GetIDByPseudo(ctx, rID)
				if err != nil {
					u.logging.Error(err)
				}

				// if id was not found - absent
				if regulationID == 0 {
					absent := entity.Absent{Pseudo: rID, ParagraphID: paragraph.ID}
					err := u.absentService.Create(ctx, absent)
					if err != nil {
						u.logging.Error(err)
					}
					continue
				}

				// get relative chapter ID
				chapterID, err := u.pseudoChapterService.GetIDByPseudo(ctx, chID)
				if err != nil {
					u.logging.Error(err)
					return err
				}

				post := fmt.Sprintf("%d#%s'>", chapterID, pID)
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

func getIDs(url string) (regID, chID, pID string) {
	fmt.Println(url)
	matchedDoc, err := regexp.MatchString(`^\/document\/cons_doc_LAW_\d+\/$`, url)
	if err != nil {
		fmt.Println(err)
		return "", "", ""
	}
	if matchedDoc {
		rID := strings.Split(strings.Split(url, "cons_doc_LAW_")[1], "/")[0]
		return rID, "", ""
	}

	matchedP, err := regexp.MatchString(`^\d+\/[a-zA-Z0-9]+\/\d+$`, url)
	if err != nil {
		fmt.Println(err)
		return "", "", ""
	}
	if matchedP {
		splited := strings.Split(url, "/")
		rID := splited[0]
		chID := splited[1]
		pID := splited[2]
		return rID, chID, pID
	}
	fmt.Println(url)
	return "", "", ""
}
