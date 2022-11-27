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
		result += space + word
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

// replace all roman with arabic numbers, if the input text is big enough split by sentences and if a sentence is huge split it by words
func CreateSpeechText(text string, bigTextLength, bigSentenceLength int) (speechText []string, err error) {
	text = dropHtml(text)
	// if the text not very big
	if len([]rune(text)) <= bigTextLength {
		// replace roman with arabic numbers
		speechText = append(speechText, replaceRomanWithArabic([]string{text})...)
		return speechText, nil
	}

	// if the text big enough split by sentences
	sentences := strings.Split(text, ". ")
	for _, sentence := range sentences {
		words := strings.Split(sentence, " ")
		if len(words) <= bigSentenceLength {
			fmt.Println(sentence)
			speechText = append(speechText, replaceRomanWithArabic([]string{sentence})...)
			continue
		}
		parts := strings.Split(sentence, ",")
		speechText = append(speechText, replaceRomanWithArabic(parts)...)
	}

	return speechText, nil
}
