package loader

import "golang.org/x/text/language"

// StringCatalog lets us store an entire catalog by tag.
type StringCatalog struct {
	Strings map[string]string
}

// Loader loads messages.
type Loader interface {
	LoadMessagesFromFile(filename string) error
	StringsByTag(tag language.Tag) (*StringCatalog, error)
}

// NewStringCatalog factory method.
func NewStringCatalog() *StringCatalog {
	return &StringCatalog{
		Strings: map[string]string{},
	}
}
