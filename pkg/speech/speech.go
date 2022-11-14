package speech

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/i-b8o/nonsense"
)

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

func CreateSpeechText(text string) (speechText []string, err error) {
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
