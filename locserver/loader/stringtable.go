package loader

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/language"
)

// StringTable loads all the languages from a base directory of locales.
type StringTable struct {
	LocalesDir string
	Matcher    *language.Matcher
	Loader     Loader
	watcher    *fsnotify.Watcher
}

// NewStringTable is a factory method for StringTable
func NewStringTable(localesDir string, watch bool, ldr Loader) (*StringTable, error) {
	strs := &StringTable{
		LocalesDir: localesDir,
		Loader:     ldr,
	}

	if watch {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil, fmt.Errorf("failed to create watcher on %s: %v", localesDir, err)
		}
		strs.watcher = watcher
	}

	return strs, nil
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

		err = st.loadMessagesFromDirectory(path.Join(st.LocalesDir, f.Name()))
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

	if st.watcher != nil {
		go st.watch()
	}

	return nil
}

func (st *StringTable) watch() {
	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-st.watcher.Events:
				log.Info().Str("name", event.Name).Uint32("op", uint32(event.Op)).Msg("Reloading strings")
				file, err := os.Open(event.Name)
				if err != nil {
					log.Error().Str("name", event.Name).Err(err).Msg("Can't get file")
					return
				}

				stat, err := file.Stat()
				if err != nil {
					log.Error().Str("name", event.Name).Err(err).Msg("Can't stat")
					return
				}

				if stat.IsDir() {
					err = st.loadMessagesFromDirectory(event.Name)
				} else {
					err = st.loadMessagesFromFile(event.Name)
				}

				if err != nil {
					log.Error().Err(err).Msg("Error reloading, not reloaded")
				}

			// watch for errors
			case err := <-st.watcher.Errors:
				log.Error().Err(err).Msg("Error from directory watcher, not reloading")
			}
		}
	}()

	<-done
}

// Close deinitializes the StringTable.
func (st *StringTable) Close() {
	st.watcher.Close()
}

func (st *StringTable) loadMessagesFromDirectory(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	st.watcher.Remove(dirname)

	var result error

	for _, f := range files {
		fullPath := path.Join(dirname, f.Name())

		if f.IsDir() {
			err = st.loadMessagesFromDirectory(fullPath) // recursive
		} else {
			err = st.loadMessagesFromFile(fullPath)
		}

		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	if st.watcher != nil {
		st.watcher.Add(dirname)
	}

	return result
}

func (st *StringTable) loadMessagesFromFile(fullPath string) error {
	var tag language.Tag
	var err error
	if st.Loader.NeedsTag() {
		dirname := filepath.Dir(fullPath)
		_, parentDir := filepath.Split(dirname)
		tagStr := parentDir

		tag, err = language.Parse(tagStr)
		if err != nil {
			return err
		}
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return st.Loader.ReadMessages(reader, &tag)
}
