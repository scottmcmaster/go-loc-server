package loader

import (
	"io"

	"golang.org/x/text/language"
)

// StringCatalog lets us store an entire catalog by tag.
type StringCatalog struct {
	Strings map[string]string
}

// Loader loads messages.
type Loader interface {
	StringsByTag(tag language.Tag) (*StringCatalog, error)

	// NeedsTag indicates whether or not this loader requires that the language tag
	// be passed or if it can be inferred from the file format.
	NeedsTag() bool

	// ReadMessages loads messages from the given reader.
	// tagStr may be ignored by the implementation if NeedsTag is false.
	ReadMessages(reader io.Reader, tagStr string) error
}

// NewStringCatalog factory method.
func NewStringCatalog() *StringCatalog {
	return &StringCatalog{
		Strings: map[string]string{},
	}
}
