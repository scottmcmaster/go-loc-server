package loader

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const header = `"Project-Id-Version: PACKAGE VERSION\n"
"Report-Msgid-Bugs-To: \n"
"POT-Creation-Date: 2013-06-04 10:20+0800\n"
"PO-Revision-Date: 2013-03-10 05:19+0800\n"
"Last-Translator: Automatically generated\n"
"Language-Team: none\n"
"Language: en_US\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
`

func TestSimplePOLoad(t *testing.T) {
	data := header + "\n" + `msgid "foo"
	msgstr "foo2"
	
	msgid "bar"
	msgstr "bar2"
	`

	reader := strings.NewReader(data)

	loader := NewPOLoader()
	err := loader.ReadMessages(reader, "en-us")

	assert.Nil(t, err)

	// Test the translation
	p := getPrinter("en-us")
	assert.Equal(t, "foo2", p.Sprintf("foo"))
	assert.Equal(t, "bar2", p.Sprintf("bar"))

	// Test the default
	p = getPrinter("en")
	assert.Equal(t, "foo", p.Sprintf("foo"))
	assert.Equal(t, "bar", p.Sprintf("bar"))
}

func TestMultiplePOLoad(t *testing.T) {
	data := header + "\n" + `msgid "foo"
	msgstr "foo2"
	
	msgid "bar"
	msgstr "bar2"
	`

	data2 := header + "\n" + `msgid "foo"
	msgstr "chinese foo"
	
	msgid "bar"
	msgstr "chinese bar"
	`

	loader := NewPOLoader()
	reader := strings.NewReader(data)
	err := loader.ReadMessages(reader, "en-us")
	assert.Nil(t, err)

	reader = strings.NewReader(data2)
	err = loader.ReadMessages(reader, "zh-cn")
	assert.Nil(t, err)

	// Test the first lang
	p := getPrinter("en-us")
	assert.Equal(t, "foo2", p.Sprintf("foo"))
	assert.Equal(t, "bar2", p.Sprintf("bar"))

	// Test the second lang
	p = getPrinter("zh-cn")
	assert.Equal(t, "chinese foo", p.Sprintf("foo"))
	assert.Equal(t, "chinese bar", p.Sprintf("bar"))
}
