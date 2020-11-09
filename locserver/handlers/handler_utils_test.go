package handlers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractLang_All(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
		URL:    &url.URL{},
	}

	req.AddCookie(&http.Cookie{Name: "lang", Value: "en-gb"})
	req.Header.Add("Accept-Language", "en-us")
	req.URL.RawQuery = "lang=zh-cn"

	lang, accept, param := ExtractLang(req)
	assert.Equal(t, lang, "en-gb")
	assert.Equal(t, accept, "en-us")
	assert.Equal(t, param, "zh-cn")
}

func TestExtractLang_None(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
		URL:    &url.URL{},
	}

	lang, accept, param := ExtractLang(req)
	assert.Equal(t, lang, "")
	assert.Equal(t, accept, "")
	assert.Equal(t, param, "")
}

func TestExtractContentType_Default(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}

	contentType := ExtractContentType(req)
	assert.Equal(t, DefaultContentType, contentType)
}

func TestExtractContentType_FromOneHeader(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	req.Header.Add("Accept", "application/json")

	contentType := ExtractContentType(req)
	assert.Equal(t, "application/json", contentType)
}

func TestExtractContentType_FromTwoHeaderFirstNotSupported(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	req.Header.Add("Accept", "text/html,application/json")

	contentType := ExtractContentType(req)
	assert.Equal(t, "application/json", contentType)
}

func TestExtractContentType_FromOneHeaderNotSupported(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	req.Header.Add("Accept", "text/html")

	contentType := ExtractContentType(req)
	assert.Equal(t, DefaultContentType, contentType)
}

func TestExtractContentType_IgnorePriority(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	req.Header.Add("Accept", "application/json;q=0.9,text/csv")

	contentType := ExtractContentType(req)
	assert.Equal(t, "application/json", contentType)
}
