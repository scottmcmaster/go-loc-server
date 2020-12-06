package loader

import (
	"io"
	"time"

	"golang.org/x/text/language"
)

// StringCatalog lets us store an entire catalog by tag.
type StringCatalog struct {
	Strings     map[string]string
	LastModTime time.Time
}

// Loader loads messages.
type Loader interface {
	StringsByTag(tag language.Tag) (*StringCatalog, error)

	// NeedsTag indicates whether or not this loader requires that the language tag
	// be passed or if it can be inferred from the file format.
	NeedsTag() bool

	// ReadMessages loads messages from the given reader.
	// tag may be ignored by the implementation if NeedsTag is false.
	ReadMessages(reader io.Reader, tag *language.Tag, modTime time.Time) error
}

// NewStringCatalog factory method.
func NewStringCatalog(modTime time.Time) *StringCatalog {
	return &StringCatalog{
		Strings:     map[string]string{},
		LastModTime: modTime,
	}
}
