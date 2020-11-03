package handlers

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

var supportedContentTypes = map[string]bool{"text/csv": true, "application/json": true}

// ExtractLang pulls the best possible language out of a request.
func ExtractLang(req *http.Request) (lang string, accept string, param string) {
	langCookie, _ := req.Cookie("lang")
	if langCookie != nil && langCookie.Name == "lang" {
		lang = langCookie.Value
	}
	accept = req.Header.Get("Accept-Language")

	langKeys, ok := req.URL.Query()["lang"]
	if ok && len(langKeys) > 0 {
		lang = langKeys[0]
	}

	return
}

// ExtractContentType gets the client's preferred response type from the request.
// Treats the order of content types in the header as the preference (i.e. it ignores preferences defined as "q" values).
func ExtractContentType(req *http.Request) (contentType string) {
	// Default to csv
	contentType = "text/csv"

	rawHeader := req.Header.Get("Accept")
	if len(rawHeader) == 0 {
		log.Debug().
			Msg("Returning default format")
		return
	}

	log.Debug().
		Str("accept", rawHeader).
		Msg("Returning requested format")

	rawTypes := strings.Split(rawHeader, ";")[0]
	types := strings.Split(rawTypes, ",")
	for _, t := range types {
		cleanType := strings.TrimSpace(t)
		if _, ok := supportedContentTypes[cleanType]; ok {
			contentType = cleanType
			return
		}
	}
	return
}
