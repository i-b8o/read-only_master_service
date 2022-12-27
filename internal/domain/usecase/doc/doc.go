package usecase_doc

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	"regexp"
	"strings"

	"github.com/i-b8o/logging"
)

type DocService interface {
	GetAll(ctx context.Context) ([]entity.Doc, error)
	Create(ctx context.Context, doc entity.Doc) (uint64, error)
	Delete(ctx context.Context, docId uint64) error
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

type PseudoDocService interface {
	Exist(ctx context.Context, pseudoID string) (bool, error)
	CreateRelationship(ctx context.Context, pseudoDoc entity.PseudoDoc) error
	DeleteRelationship(ctx context.Context, docID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type PseudoChapterService interface {
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type docUsecase struct {
	docService           DocService
	chapterService       ChapterService
	paragraphService     ParagraphService
	absentService        AbsentService
	pseudoDocService     PseudoDocService
	pseudoChapterService PseudoChapterService
	logging              logging.Logger
}

func NewDocUsecase(docService DocService, chapterService ChapterService, paragraphService ParagraphService, absentService AbsentService, pseudoDocService PseudoDocService, pseudoChapterService PseudoChapterService, logging logging.Logger) *docUsecase {
	return &docUsecase{docService: docService, chapterService: chapterService, paragraphService: paragraphService, absentService: absentService, pseudoDocService: pseudoDocService, pseudoChapterService: pseudoChapterService, logging: logging}
}

func (u docUsecase) Exist(ctx context.Context, pseudo string) (bool, error) {
	return u.pseudoDocService.Exist(ctx, pseudo)
}

func (u docUsecase) GetAll(ctx context.Context) ([]entity.Doc, error) {
	return u.docService.GetAll(ctx)
}

func (u docUsecase) CreateDoc(ctx context.Context, doc entity.Doc) (uint64, error) {
	// create a doc
	ID, err := u.docService.Create(ctx, doc)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// create an id-pseudoId relationship
	err = u.pseudoDocService.CreateRelationship(ctx, entity.PseudoDoc{ID: ID, PseudoId: doc.Pseudo})
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// mark the doc as done
	err = u.absentService.Done(ctx, doc.Pseudo)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	return ID, nil
}

func (u docUsecase) GetAbsents(ctx context.Context) ([]*entity.Absent, error) {
	absents, err := u.absentService.GetAll(ctx)
	if err != nil {
		u.logging.Error(err)
		return nil, err
	}

	return absents, nil
}

func (u docUsecase) DeleteDoc(ctx context.Context, ID uint64) error {
	// delete a doc
	err := u.docService.Delete(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	// delete the id-pseudoId relationship
	err = u.pseudoDocService.DeleteRelationship(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	return nil
}

func (u docUsecase) GenerateLinks(ctx context.Context, docID uint64) error {
	// get IDs for every chapter in the doc
	chIDs, err := u.chapterService.GetAllIds(ctx, docID)
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
					err := u.absentService.Create(ctx, absent)
					if err != nil {
						u.logging.Error(err)
					}
					continue
				}
				// link for an entire document.
				if chID == "" {

					// get relative doc ID
					docID, err := u.pseudoDocService.GetIDByPseudo(ctx, rID)
					if err != nil {
						u.logging.Error(err)
					}
					if docID == 0 {
						absent := entity.Absent{Pseudo: rID, ParagraphID: paragraph.ID}
						err := u.absentService.Create(ctx, absent)
						if err != nil {
							u.logging.Error(err)
						}
						continue
					}
					post := fmt.Sprintf("%d'>", docID)
					content = strings.Replace(content, aLink, "<a href='"+post, 1)
				}

				// link for a paragraph
				// get relative doc ID
				docID, err := u.pseudoDocService.GetIDByPseudo(ctx, rID)
				if err != nil {
					u.logging.Error(err)
				}

				// if id was not found - absent
				if docID == 0 {
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
	return "", "", ""
}