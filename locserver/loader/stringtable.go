package loader

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"

	"golang.org/x/text/language"
)

// StringTable loads all the languages from a base directory of locales.
type StringTable struct {
	LocalesDir string
	Matcher    *language.Matcher
	Loader     Loader
}

// Load loads the languages from the configured local directory.
func (st *StringTable) Load() error {
	log.Info().Str("localesdir", st.LocalesDir).Msg("Loading locales")

	files, err := ioutil.ReadDir(st.LocalesDir)
	if err != nil {
		return err
	}

	tags := []language.Tag{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		t, err := language.Parse(f.Name())
		if err != nil {
			log.Warn().Err(err).Str("locale", f.Name()).Msg("Unable to parse locale directory name")
			continue
		}
		tags = append(tags, t)

		err = st.loadMessageFromDirectory(path.Join(st.LocalesDir, f.Name()))
		if err != nil {
			log.Warn().Err(err).Str("locale", f.Name()).Msg("Error reading locale directory")
		}
	}

	if len(tags) == 0 {
		return errors.New("no language tags found")
	}

	log.Info().Interface("tags", tags).Msg("Creating matcher")
	matcher := language.NewMatcher(tags)
	st.Matcher = &matcher

	return nil
}

func (st *StringTable) loadMessageFromDirectory(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	var result error

	for _, f := range files {
		fullPath := path.Join(dirname, f.Name())

		if f.IsDir() {
			err = st.loadMessageFromDirectory(fullPath)
		} else {
			tagStr := ""
			if st.Loader.NeedsTag() {
				_, parentDir := filepath.Split(dirname)
				tagStr = parentDir
			}
			err = st.Loader.LoadMessagesFromFile(fullPath, tagStr)
		}

		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
