package usecase_regulation

import (
	"context"
	"fmt"
	"prod_serv/internal/domain/entity"
	"regexp"
	"strconv"
	"strings"

	"github.com/i-b8o/nonsense"
)

type RegulationService interface {
	GetOne(ctx context.Context, regulationID uint64) (entity.Regulation, error)
	GetAll(ctx context.Context) ([]entity.Regulation, error)
	Create(ctx context.Context, regulation entity.Regulation) (string, error)
	DeleteRegulation(ctx context.Context, regulationID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}
type ChapterService interface {
	GetAllById(ctx context.Context, regulationID uint64) ([]entity.Chapter, error)
	GetOrderNum(ctx context.Context, id uint64) (orderNum uint64, err error)
	DeleteForRegulation(ctx context.Context, regulationID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}
type ParagraphService interface {
	GetAllById(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error)
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]entity.Paragraph, error)
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type AbsentService interface {
	Create(ctx context.Context, absent entity.Absent) error
	Done(ctx context.Context, pseudo string) error
}

type LinkService interface {
	GetAll(ctx context.Context) ([]*entity.Link, error)
	GetAllByChapterID(ctx context.Context, chapterID uint64) ([]*entity.Link, error)
	// Create(ctx context.Context, link entity.Link) error
	GetOneByParagraphID(ctx context.Context, paragraphID, regregulationID uint64) (entity.Link, error)
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type SpeechService interface {
	GetAllById(ctx context.Context, paragraphID uint64) ([]entity.Speech, error)
	DeleteForParagraph(ctx context.Context, paragraphID uint64) error
}

type regulationUsecase struct {
	regulationService RegulationService
	chapterService    ChapterService
	paragraphService  ParagraphService
	linkService       LinkService
	speechService     SpeechService
	absentService     AbsentService
}

func NewRegulationUsecase(regulationService RegulationService, chapterService ChapterService, paragraphService ParagraphService, linkService LinkService, speechService SpeechService, absentService AbsentService) *regulationUsecase {
	return &regulationUsecase{regulationService: regulationService, chapterService: chapterService, paragraphService: paragraphService, linkService: linkService, speechService: speechService, absentService: absentService}
}

func (u regulationUsecase) CreateRegulation(ctx context.Context, regulation entity.Regulation) string {
	id, err := u.regulationService.Create(ctx, regulation)
	if err != nil {
		return ""
	}
	u.absentService.Done(ctx, regulation.Pseudo)
	if err != nil {
		return ""
	}

	return id
}

func (u regulationUsecase) GetFullRegulationByID(ctx context.Context, regulationID uint64) (entity.Regulation, error) {
	regulation, err := u.regulationService.GetOne(ctx, regulationID)
	if err != nil {
		return entity.Regulation{}, err
	}
	chapters, err := u.chapterService.GetAllById(ctx, regulationID)
	if err != nil {
		return entity.Regulation{}, err
	}

	for _, chapter := range chapters {
		paragraphs, err := u.paragraphService.GetAllById(ctx, chapter.ID)
		if err != nil {
			return entity.Regulation{}, err
		}
		chapter.Paragraphs = paragraphs
	}

	regulation.Chapters = chapters

	return regulation, nil
}

func (u regulationUsecase) GetDartFullRegulationByID(ctx context.Context, regulationID uint64) string {
	regulation, err := u.regulationService.GetOne(ctx, regulationID)
	if err != nil {
		return ""
	}
	chapters, err := u.chapterService.GetAllById(ctx, regulationID)
	if err != nil {
		return ""
	}

	dartClass := `
	import 'paragraph.dart';
	import 'chapter.dart';
	
	class Regulation {
		static const int id = %d;
		static const String name = "%s";
		static const String abbreviation = "%s";
		static const List<Chapter> chapters = <Chapter>[
			%s
		];
	}
	`

	chaptersDartString, _ := u.chaptersDart(ctx, chapters)
	return fmt.Sprintf(dartClass, regulationID, regulation.Name, regulation.Abbreviation, chaptersDartString)
}

func (u regulationUsecase) chaptersDart(ctx context.Context, chapters []entity.Chapter) (dartChaptersString string, err error) {
	dartChapter := `Chapter(id: %d, name: "%s", num: "%s", orderNum: %d , paragraphs: [
		%s
	]),`
	for _, chapter := range chapters {
		paragraphs, err := u.paragraphService.GetAllById(ctx, chapter.ID)
		if err != nil {
			return dartChapter, err
		}
		dartPar := paragraphsDart(ctx, paragraphs, u)

		num := ""
		if len(chapter.Num) > 0 {
			num = chapter.Num
		}
		name := strings.Replace(chapter.Name, "\n", "", -1)
		temp := fmt.Sprintf(dartChapter, chapter.ID, name, num, chapter.OrderNum, dartPar)
		dartChaptersString += temp
	}
	return dartChaptersString, nil
}

func speachText(contentSlice []string) string {
	start := `[`
	end := `]`

	var result string

	for _, part := range contentSlice {
		part = strings.ReplaceAll(part, "'", `"`)
		str := fmt.Sprintf(`'%s',`, part)
		result += str
	}
	result = result[:len(result)-1]
	re := regexp.MustCompile(`\r?\n`)
	result = re.ReplaceAllString(result, " ")
	return start + result + end
}

func paragraphsDart(ctx context.Context, paragraphs []entity.Paragraph, u regulationUsecase) (dartParagraphsList string) {
	for _, p := range paragraphs {
		text := strings.Replace(p.Content, "\n", "", -1)
		text = strings.ReplaceAll(text, `'`, `"`)

		var speechTextSlice []string
		if p.ID > 0 {
			speechSlice, err := u.speechService.GetAllById(ctx, p.ID)
			if err != nil {
				return ""
			}
			for _, t := range speechSlice {
				speechTextSlice = append(speechTextSlice, t.Content)
			}
		} else {
			speechTextSlice = append(speechTextSlice, p.Content)
		}

		textToSpeech := speachText(speechTextSlice)
		dartParagraphsList += fmt.Sprintf(`		Paragraph(id: %d, num: %d, isTable: %t,isNFT: %t, paragraphClass: "%s", content: '%s', chapterID: %d, textToSpeech: %s),
		`, p.ID, p.Num, p.IsTable, p.IsNFT, p.Class, text, p.ChapterID, textToSpeech)
	}
	return dartParagraphsList
}

func (u regulationUsecase) linksDart(ctx context.Context, links []*entity.Link) (dartLinksList string) {

	for _, l := range links {
		num, err := u.chapterService.GetOrderNum(ctx, l.ChapterID)
		if err != nil {
			fmt.Println(err)
		}
		dartLinksList += fmt.Sprintf(`		Link(id: %d, chapterNum: %d, paragraphNum: %d, rid: %d),
		`, l.ID, num, l.ParagraphNum, l.RID)
	}
	return dartLinksList
}

func (u regulationUsecase) AllLinksDart(ctx context.Context, regulationID uint64) string {
	chapters, err := u.chapterService.GetAllById(ctx, regulationID)
	if err != nil {
		return ""
	}
	var links []*entity.Link

	for _, chapter := range chapters {
		l, _ := u.linkService.GetAllByChapterID(ctx, chapter.ID)
		links = append(links, l...)
	}

	for _, l := range links {
		l.RID = regulationID
	}

	dartClass := `
import 'link.dart';
	
class AllLinks {
	static const List<Link> links = <Link>[
		%s
	];
}
	`

	linkssDartString := u.linksDart(ctx, links)

	return fmt.Sprintf(dartClass, linkssDartString)
}

func (u regulationUsecase) GetDocumentRoot(ctx context.Context, stringID string) (entity.Regulation, []entity.Chapter) {
	uint64ID, err := strconv.ParseUint(stringID, 10, 64)
	if err != nil {
		return entity.Regulation{}, nil
	}
	regulation, err := u.regulationService.GetOne(ctx, uint64ID)
	if err != nil {
		return entity.Regulation{}, nil
	}

	regulation.Name = nonsense.Capitalize(regulation.Name)
	chapters, err := u.chapterService.GetAllById(ctx, uint64ID)
	if err != nil {
		return entity.Regulation{}, nil
	}
	return regulation, chapters
}

func (u regulationUsecase) GetDocuments(ctx context.Context) []entity.Regulation {
	regulations, err := u.regulationService.GetAll(ctx)
	if err != nil {
		return nil
	}

	return regulations
}

func (u regulationUsecase) GenerateLinks(ctx context.Context, regulationID uint64) error {
	chapters, err := u.chapterService.GetAllById(ctx, regulationID)
	if err != nil {
		return err
	}

	for _, chapter := range chapters {
		paragraphs, err := u.paragraphService.GetWithHrefs(ctx, chapter.ID)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		for _, paragraph := range paragraphs {
			content := paragraph.Content
			re := regexp.MustCompile("<a href='.+'>")
			links := re.FindAllString(content, -1)
			// <a href='357933/f03432cdbb8d43bc45fb191c0e0e91393229c429/335782'>
			for _, aLink := range links {
				hrefRaw := strings.Split(aLink, "<a href='")[1]
				href := strings.Split(hrefRaw, "'>")[0]
				// Link for document
				matched, err := regexp.MatchString(`^\/document\/cons_doc_LAW_\d+\/$`, href)
				if err != nil {
					fmt.Println(err.Error())
				}

				if matched {
					fmt.Printf("This: %s", href)
					p := strings.Split(href, "/document/cons_doc_LAW_")[1]
					ID := strings.Split(p, "/")[0]
					rID, err := u.regulationService.GetIDByPseudo(ctx, ID)
					if err != nil {
						fmt.Printf("href = %s,ID: %s, error: %s", href, ID, err.Error())
					}
					if rID == 0 {
						fmt.Printf("ID: %d, %s, %d|\n", rID, ID, paragraph.ID)
						absent := entity.Absent{Pseudo: ID, ParagraphID: paragraph.ID}
						err := u.absentService.Create(ctx, absent)
						if err != nil {
							fmt.Println(err.Error())
						}
						continue
					}
				}

				IDs := strings.Split(href, "/")
				if (len(IDs) < 3) || (len(IDs[0]) == 0) || (len(IDs[1]) == 0) || (len(IDs[2]) == 0) {
					absent := entity.Absent{Pseudo: href, ParagraphID: paragraph.ID}
					u.absentService.Create(ctx, absent)
					continue
				}
				regID, err := u.regulationService.GetIDByPseudo(ctx, IDs[0])
				if err != nil {
					fmt.Printf("href = %s,ID: %s, error: %s", href, IDs[0], err.Error())
				}
				if regID == 0 {
					absent := entity.Absent{Pseudo: IDs[0], ParagraphID: paragraph.ID}
					u.absentService.Create(ctx, absent)
					continue
				}
				chID, err := u.chapterService.GetIDByPseudo(ctx, IDs[1])

				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				pID := IDs[2]

				post := fmt.Sprintf("%d/%d/%s>", regID, chID, pID)
				// post := fmt.Sprintf("#%s'>", href)
				content = strings.Replace(content, aLink, "<a href='"+post, 1)
			}
			err := u.paragraphService.UpdateOne(ctx, content, paragraph.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u regulationUsecase) DeleteRegulation(ctx context.Context, regulationID uint64) error {
	chapters, err := u.chapterService.GetAllById(ctx, regulationID)
	if err != nil {
		return err
	}

	for _, chapter := range chapters {
		err = u.linkService.DeleteForChapter(ctx, chapter.ID)
		if err != nil {
			return err
		}
		paragraphs, err := u.paragraphService.GetAllById(ctx, chapter.ID)
		if err != nil {
			return err
		}
		for _, paragraph := range paragraphs {
			err = u.speechService.DeleteForParagraph(ctx, paragraph.ID)
			if err != nil {
				return err
			}
		}
		err = u.paragraphService.DeleteForChapter(ctx, chapter.ID)
		if err != nil {
			return err
		}

	}
	err = u.chapterService.DeleteForRegulation(ctx, regulationID)
	if err != nil {
		return err
	}
	err = u.regulationService.DeleteRegulation(ctx, regulationID)
	if err != nil {
		return err
	}

	return nil
}
