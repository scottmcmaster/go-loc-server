package loader

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type langmessage struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	Translation string `json:"translation"`
}

type langmessages struct {
	Language string        `json:"language"`
	Messages []langmessage `json:"messages"`
}

// GoTextJSONLoader loads strings from files in the JSON format supported by gotext.
type GoTextJSONLoader struct {
	catalogsByTagStr map[string]*StringCatalog
}

// NewGoTextJSONLoader factory method.
func NewGoTextJSONLoader() *GoTextJSONLoader {
	return &GoTextJSONLoader{
		catalogsByTagStr: map[string]*StringCatalog{},
	}
}

// StringsByTag gets the string table for the given language tag.
func (ldr *GoTextJSONLoader) StringsByTag(tag language.Tag) (*StringCatalog, error) {
	if cat, ok := ldr.catalogsByTagStr[tag.String()]; ok {
		return cat, nil
	}
	return nil, errors.New("catalog not found for tag " + tag.String())
}

// LoadMessagesFromFile implements the Loader interface.
func (ldr *GoTextJSONLoader) LoadMessagesFromFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lm := langmessages{}

	err = json.Unmarshal([]byte(file), &lm)
	if err != nil {
		return err
	}

	t, err := language.Parse(lm.Language)
	if err != nil {
		return err
	}

	tagStr := t.String()
	ldr.catalogsByTagStr[tagStr] = NewStringCatalog()

	for _, m := range lm.Messages {
		log.Debug().Str("languagetag", tagStr).
			Str("id", m.ID).
			Str("translation", m.Translation).
			Msg("Loading string")
		message.SetString(t, m.ID, m.Translation)

		ldr.catalogsByTagStr[tagStr].Strings[m.ID] = m.Translation
	}

	return nil
}