package loader

import "golang.org/x/text/language"

// StringCatalog lets us store an entire catalog by tag.
type StringCatalog struct {
	Strings map[string]string
}

// Loader loads messages.
type Loader interface {
	// LoadMessagesFromFile loads messages from the given file.
	// tagStr may be ignored by the implementation if NeedsTag is false.
	LoadMessagesFromFile(filename string, tagStr string) error

	StringsByTag(tag language.Tag) (*StringCatalog, error)

	// NeedsTag indicates whether or not this loader requires that the language tag
	// be passed or if it can be inferred from the file format.
	NeedsTag() bool
}

// NewStringCatalog factory method.
func NewStringCatalog() *StringCatalog {
	return &StringCatalog{
		Strings: map[string]string{},
	}
}
