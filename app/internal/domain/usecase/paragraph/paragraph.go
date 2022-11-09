package usecase_paragraph

import (
	"context"
	"fmt"
	"prod_serv/internal/domain/entity"
	"regexp"
	"strconv"
	"strings"

	"github.com/i-b8o/nonsense"
)

type ParagraphService interface {
	CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error
}

type ChapterService interface {
	GetOneById(ctx context.Context, chapterID uint64) (entity.Chapter, error)
}

type LinkService interface {
	Create(ctx context.Context, link entity.Link) error
}

type SpeechService interface {
	Create(ctx context.Context, speech entity.Speech) (string, error)
}

type paragraphUsecase struct {
	paragraphService ParagraphService
	chapterService   ChapterService
	linkService      LinkService
	speechService    SpeechService
}

func NewParagraphUsecase(paragraphService ParagraphService, chapterService ChapterService, linkService LinkService, speechService SpeechService) *paragraphUsecase {
	return &paragraphUsecase{paragraphService: paragraphService, chapterService: chapterService, linkService: linkService, speechService: speechService}
}

func (u paragraphUsecase) CreateParagraphs(ctx context.Context, paragraphs []entity.Paragraph) error {
	if len(paragraphs) == 0 {
		return nil
	}
	ch, err := u.chapterService.GetOneById(ctx, paragraphs[0].ChapterID)
	if err != nil {
		return err
	}
	for _, p := range paragraphs {
		if p.ID > 0 {
			u.linkService.Create(ctx, entity.Link{ID: p.ID, ParagraphNum: p.Num, ChapterID: p.ChapterID, RID: ch.RegulationID})
			speechTextSlice, err := createSpeechText(p.Content)
			if err != nil {
				return err
			}
			for i, text := range speechTextSlice {
				speech := entity.Speech{ParagraphID: p.ID, Content: text, OrderNum: uint64(i)}
				_, err := u.speechService.Create(ctx, speech)
				if err != nil {
					return err
				}
			}
		}

		// When a paragraph has IDs inside
		hasIDsInside := strings.Contains(p.Content, "<a id=")
		if hasIDsInside {
			re := regexp.MustCompile(`<a id='(.*?)'`)
			matches := re.FindAllString(p.Content, -1)
			for _, match := range matches {
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
				subIndexStr := re.FindString(match)
				subIndexUint64, err := strconv.ParseUint(subIndexStr, 10, 64)
				if err != nil {
					return err
				}
				u.linkService.Create(ctx, entity.Link{ID: subIndexUint64, ParagraphNum: p.Num, ChapterID: p.ChapterID, RID: ch.RegulationID})
			}
		}
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

func dropHtml(s string) string {
	const regex = `<.*?>`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}

func replaceRomanWithArabicString(text string) (result string) {
	words := strings.Split(text, " ")
	for i, word := range words {
		space := ""
		if i > 0 {
			space = " "
		}
		if nonsense.IsRoman(word) {
			arabic, _ := nonsense.ToIndoArabic(word)
			result += fmt.Sprintf("%s%d", space, arabic)
			continue
		}
		result += " " + word
	}

	return result
}

func replaceRomanWithArabic(text []string) []string {
	var result []string
	for _, str := range text {
		result = append(result, replaceRomanWithArabicString(str))
	}
	return result
}

func createSpeechText(text string) (speechText []string, err error) {
	text = dropHtml(text)
	if len([]rune(text)) <= 250 {
		speechText = append(speechText, replaceRomanWithArabic([]string{text})...)
		return speechText, nil
	}

	sentences := strings.Split(text, ". ")
	for _, sentence := range sentences {
		words := strings.Split(sentence, " ")
		if len(words) <= 40 {
			speechText = append(speechText, replaceRomanWithArabic([]string{sentence})...)
			// fmt.Println("here " + speechText)
			continue
		}
		parts := strings.Split(sentence, ",")
		speechText = append(speechText, replaceRomanWithArabic(parts)...)
	}

	return speechText, nil
}
