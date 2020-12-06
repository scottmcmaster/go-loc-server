package loader

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// xliff is a stripped-down representation of the full XLIFF 2.0 schema.
type xliff struct {
	XMLName xml.Name `xml:"xliff"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Version string   `xml:"version,attr"`
	SrcLang string   `xml:"srcLang,attr"`
	TrgLang string   `xml:"trgLang,attr"`
	File    struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
		Unit []struct {
			Text    string `xml:",chardata"`
			Segment []struct {
				Text   string `xml:",chardata"`
				ID     string `xml:"id,attr"`
				Source string `xml:"source"`
				Target string `xml:"target"`
			} `xml:"segment"`
		} `xml:"unit"`
	} `xml:"file"`
}

// XLIFF2Loader loads strings from files in the XLIFF 2 format.
type XLIFF2Loader struct {
	catalogsByTagStr map[string]*StringCatalog
}

// NewXLIFF2Loader factory method.
func NewXLIFF2Loader() *XLIFF2Loader {
	return &XLIFF2Loader{
		catalogsByTagStr: map[string]*StringCatalog{},
	}
}

// StringsByTag gets the string table for the given language tag.
func (ldr *XLIFF2Loader) StringsByTag(tag language.Tag) (*StringCatalog, error) {
	if cat, ok := ldr.catalogsByTagStr[tag.String()]; ok {
		return cat, nil
	}
	return nil, errors.New("catalog not found for tag " + tag.String())
}

// NeedsTag implements the Loader interface.
func (ldr *XLIFF2Loader) NeedsTag() bool {
	// Not needed because the langauge is embedded in the file.
	return false
}

// ReadMessages implements the Loader interface.
func (ldr *XLIFF2Loader) ReadMessages(reader io.Reader, tag *language.Tag, modTime time.Time) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var xlf xliff
	xml.Unmarshal(data, &xlf)

	t, err := language.Parse(xlf.TrgLang)
	if err != nil {
		return err
	}

	tagStr := t.String()
	ldr.catalogsByTagStr[tagStr] = NewStringCatalog(modTime)

	for _, u := range xlf.File.Unit {
		for _, seg := range u.Segment {
			log.Debug().Str("languagetag", tagStr).
				Str("id", seg.ID).
				Str("translation", seg.Target).
				Msg("Loading string")
			message.SetString(t, seg.ID, seg.Target)

			ldr.catalogsByTagStr[tagStr].Strings[seg.ID] = seg.Target
		}
	}

	return nil
}
