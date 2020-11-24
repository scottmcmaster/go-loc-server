package loader

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const pattern = `msgid "(.+)"\nmsgstr "(.+)"`

// POLoader loads strings from files in the gettext PO format.
type POLoader struct {
	catalogsByTagStr map[string]*StringCatalog
}

// NewPOLoader factory method.
func NewPOLoader() *POLoader {
	return &POLoader{
		catalogsByTagStr: map[string]*StringCatalog{},
	}
}

// StringsByTag gets the string table for the given language tag.
func (ldr *POLoader) StringsByTag(tag language.Tag) (*StringCatalog, error) {
	if cat, ok := ldr.catalogsByTagStr[tag.String()]; ok {
		return cat, nil
	}
	return nil, errors.New("catalog not found for tag " + tag.String())
}

// NeedsTag implements the Loader interface.
func (ldr *POLoader) NeedsTag() bool {
	// Needed because the langauge is not embedded in the file.
	return true
}

// ReadMessages implements the Loader interface.
func (ldr *POLoader) ReadMessages(reader io.Reader, tagStr string) error {

	return nil
}

// LoadMessagesFromFile implements the Loader interface.
func (ldr *POLoader) LoadMessagesFromFile(filename string, tagStr string) error {
	if len(tagStr) == 0 {
		return errors.New("tag string is required by PO loader")
	}

	t, err := language.Parse(tagStr)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	ldr.catalogsByTagStr[tagStr] = NewStringCatalog()

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(buf), -1)

	pairs := map[string]string{}
	for _, array := range matches {
		log.Debug().Str("languagetag", tagStr).
			Str("id", array[1]).
			Str("translation", array[2]).
			Msg("Loading string")
		message.SetString(t, array[1], array[2])

		ldr.catalogsByTagStr[tagStr].Strings[array[1]] = array[2]
		pairs[array[1]] = array[2]
	}

	return nil
}
