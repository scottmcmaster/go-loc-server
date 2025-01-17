package loader

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleGoTextJSONLoad(t *testing.T) {
	data := `{
		"language": "en-us",
		"messages": [{
			"id": "foo",
			"message": "foo",
			"translation": "foo2"
		  },
		  {
			"id": "bar",
			"message": "bar",
			"translation": "bar2"
		  }
		]
	  }`

	reader := strings.NewReader(data)

	loader := NewGoTextJSONLoader()
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

func TestMultipleGoTextJSONLoad(t *testing.T) {
	data := `{
		"language": "en-us",
		"messages": [{
			"id": "foo",
			"message": "foo",
			"translation": "foo2"
		  },
		  {
			"id": "bar",
			"message": "bar",
			"translation": "bar2"
		  }
		]
	  }`

	data2 := `{
		"language": "zh-cn",
		"messages": [{
			"id": "foo",
			"message": "foo",
			"translation": "chinese foo"
		  },
		  {
			"id": "bar",
			"message": "bar",
			"translation": "chinese bar"
		  }
		]
	  }`

	loader := NewGoTextJSONLoader()
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
