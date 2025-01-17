package loader

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"time"

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
func (ldr *POLoader) ReadMessages(reader io.Reader, tag *language.Tag, modTime time.Time) error {
	if tag == nil {
		return errors.New("tag string is required by PO loader")
	}

	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	tagStr := tag.String()
	ldr.catalogsByTagStr[tagStr] = NewStringCatalog(modTime)

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(buf), -1)

	pairs := map[string]string{}
	for _, array := range matches {
		log.Debug().Str("languagetag", tagStr).
			Str("id", array[1]).
			Str("translation", array[2]).
			Msg("Loading string")
		message.SetString(*tag, array[1], array[2])

		ldr.catalogsByTagStr[tagStr].Strings[array[1]] = array[2]
		pairs[array[1]] = array[2]
	}

	return nil
}
