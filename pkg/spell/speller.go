package spell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	ErrCheckingText     = "error checking text"
	ErrDecodingResponse = "error decoding response"
)

type Speller interface {
	CheckText(text string) ([]SpellingError, error)
	FormatErrors(errors []SpellingError) string
}

type YandexSpeller struct {
	SpellerURL string
}

func NewYandexSpeller(yandexSpellerURL string) *YandexSpeller {
	return &YandexSpeller{
		SpellerURL: yandexSpellerURL,
	}
}

type SpellingError struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func (s *YandexSpeller) CheckText(text string) ([]SpellingError, error) {
	resp, err := http.PostForm(s.SpellerURL, url.Values{"text": {text}})
	if err != nil {
		return nil, fmt.Errorf(ErrCheckingText+": %v", err)
	}
	defer resp.Body.Close()

	var spellingErrors []SpellingError
	if err = json.NewDecoder(resp.Body).Decode(&spellingErrors); err != nil {
		return nil, fmt.Errorf(ErrDecodingResponse+": %v", err)
	}

	return spellingErrors, nil
}

func (s *YandexSpeller) FormatErrors(errors []SpellingError) string {
	var builder strings.Builder
	for _, e := range errors {
		builder.WriteString(fmt.Sprintf("Word '%s' at position %d is misspelled. Suggestions: %s. ", e.Word, e.Pos, strings.Join(e.S, ", ")))
	}
	return builder.String()
}
