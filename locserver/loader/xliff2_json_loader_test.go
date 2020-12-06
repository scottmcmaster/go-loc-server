package loader

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleXLIFF2Load(t *testing.T) {
	data := `<xliff xmlns="urn:oasis:names:tc:xliff:document:2.0" version="2.0" srcLang="en-us" trgLang="en-us">
	<file id="en-us">
	 <unit>
	  <segment id="foo">
	   <source>foo/source>
	   <target>foo2</target>
	  </segment>
	  <segment id="bar">
	   <source>bar</source>
	   <target>bar2</target>
	  </segment>
	 </unit>
	</file>
   </xliff>`

	reader := strings.NewReader(data)

	loader := NewXLIFF2Loader()
	err := loader.ReadMessages(reader, nil, time.Now())

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

func TestMultipleXLIFF2Load(t *testing.T) {
	data := `<xliff xmlns="urn:oasis:names:tc:xliff:document:2.0" version="2.0" srcLang="en-us" trgLang="en-us">
	<file id="en-us">
	 <unit>
	  <segment id="foo">
	   <source>foo/source>
	   <target>foo2</target>
	  </segment>
	  <segment id="bar">
	   <source>bar</source>
	   <target>bar2</target>
	  </segment>
	 </unit>
	</file>
   </xliff>`

	data2 := `<xliff xmlns="urn:oasis:names:tc:xliff:document:2.0" version="2.0" srcLang="en-us" trgLang="en-us">
   <file id="en-us">
	<unit>
	 <segment id="foo">
	  <source>foo/source>
	  <target>chinese foo</target>
	 </segment>
	 <segment id="bar">
	  <source>bar</source>
	  <target>chinese bar</target>
	 </segment>
	</unit>
   </file>
  </xliff>`

	loader := NewXLIFF2Loader()
	reader := strings.NewReader(data)
	err := loader.ReadMessages(reader, nil, time.Now())
	assert.Nil(t, err)

	reader = strings.NewReader(data2)
	err = loader.ReadMessages(reader, nil, time.Now())
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
