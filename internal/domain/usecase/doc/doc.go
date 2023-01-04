package usecase_doc

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	"regexp"
	"strings"
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
	UpdateOne(ctx context.Context, id, chapterID uint64, content string) error
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
}

func NewDocUsecase(docService DocService, chapterService ChapterService, paragraphService ParagraphService, absentService AbsentService, pseudoDocService PseudoDocService, pseudoChapterService PseudoChapterService) *docUsecase {
	return &docUsecase{docService: docService, chapterService: chapterService, paragraphService: paragraphService, absentService: absentService, pseudoDocService: pseudoDocService, pseudoChapterService: pseudoChapterService}
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
		return 0, err
	}
	// create an id-pseudoId relationship
	err = u.pseudoDocService.CreateRelationship(ctx, entity.PseudoDoc{ID: ID, PseudoId: doc.Pseudo})
	if err != nil {
		return 0, err
	}

	// mark the doc as done
	err = u.absentService.Done(ctx, doc.Pseudo)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (u docUsecase) GetAbsents(ctx context.Context) ([]*entity.Absent, error) {
	absents, err := u.absentService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return absents, nil
}

func (u docUsecase) DeleteDoc(ctx context.Context, ID uint64) error {
	// delete a doc
	err := u.docService.Delete(ctx, ID)
	if err != nil {
		return err
	}
	// delete the id-pseudoId relationship
	err = u.pseudoDocService.DeleteRelationship(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

func (u docUsecase) GenerateLinks(ctx context.Context, docID uint64) error {
	// get IDs for every chapter in the doc
	chIDs, err := u.chapterService.GetAllIds(ctx, docID)
	if err != nil {
		return err
	}

	for _, chId := range chIDs {
		// get only paragraphs with links inside
		paragraphs, err := u.paragraphService.GetParagraphsWithHrefs(ctx, chId)
		if err != nil {
			return err
		}

		for _, paragraph := range paragraphs {

			// get links
			content := paragraph.Content
			re := regexp.MustCompile("<a href='.+?'>")
			links := re.FindAllString(content, -1)

			for _, aLink := range links {
				fmt.Println("alink|" + aLink + "|")
				// <a href='/document/cons_doc_LAW_435882/'>
				hrefRaw := strings.Split(aLink, "<a href='")[1]
				fmt.Println("hrefRaw:" + hrefRaw + "|")
				href := strings.Split(hrefRaw, "'>")[0]
				fmt.Println("href:" + href + "|")
				// get if exist IDs for reggulation, chapter and paragrap
				rID, chID, pID := getIDs(href)
				fmt.Printf("first rID: %s, chID: %s, pID: %s", rID, chID, pID)
				// something wrong with the link - absent
				fmt.Println("0")
				if rID == "" {
					fmt.Println("1")
					absent := entity.Absent{Pseudo: href, ChapterID: paragraph.ChapterID, ParagraphID: paragraph.ID}
					err := u.absentService.Create(ctx, absent)
					if err != nil {
						fmt.Println("AAAAAAAAAAAAAAAAAA1")
						return err
					}
					fmt.Println("2")
					continue
				}
				// link for an entire document.
				if chID == "" {
					fmt.Println("3")
					// get relative doc ID
					docID, _ := u.pseudoDocService.GetIDByPseudo(ctx, rID)
					fmt.Printf("docID: %d\n", docID)
					if docID == 0 {
						fmt.Println("4")
						absent := entity.Absent{Pseudo: rID, ChapterID: paragraph.ChapterID, ParagraphID: paragraph.ID}
						err := u.absentService.Create(ctx, absent)
						if err != nil {
							fmt.Println("AAAAAAAAAAAAAAAAAA2")
							return err
						}
						fmt.Println("5")
						continue
					}
					fmt.Println("6")
					post := fmt.Sprintf("%d'>", docID)
					content = strings.Replace(content, aLink, "<a href='/doc/"+post, 1)
					fmt.Printf("1 post:%s, content: %s", post, content)
				} else {
					fmt.Println("7")
					// link for a paragraph
					// get relative doc ID
					docID, _ := u.pseudoDocService.GetIDByPseudo(ctx, rID)

					// if id was not found - absent
					if docID == 0 {
						absent := entity.Absent{Pseudo: rID, ChapterID: paragraph.ChapterID, ParagraphID: paragraph.ID}
						err := u.absentService.Create(ctx, absent)
						if err != nil {
							fmt.Println("AAAAAAAAAAAAAAAAAA3")
							return err
						}
						continue
					}

					// get relative chapter ID
					chapterID, err := u.pseudoChapterService.GetIDByPseudo(ctx, chID)
					if err != nil {
						fmt.Println("AAAAAAAAAAAAAAAAAA4")
						return err
					}
					fmt.Printf("rID: %s, chID: %s, pID: %s", rID, chID, pID)
					post := fmt.Sprintf("%d#%s'>", chapterID, pID)
					content = strings.Replace(content, aLink, "<a href='"+post, 1)
				}
			}

			fmt.Println("2 content|" + content + "|")
			err := u.paragraphService.UpdateOne(ctx, paragraph.ID, paragraph.ChapterID, content)
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
