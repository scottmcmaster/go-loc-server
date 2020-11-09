package handlers

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// DefaultContentType is the output content type if no other is found on the request
const DefaultContentType = "text/csv"

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
		param = langKeys[0]
	}

	return
}

// ExtractContentType gets the client's preferred response type from the request.
// Treats the order of content types in the header as the preference (i.e. it ignores preferences defined as "q" values).
func ExtractContentType(req *http.Request) (contentType string) {
	// Default to csv
	contentType = DefaultContentType

	rawHeader := req.Header.Get("Accept")
	if len(rawHeader) == 0 {
		log.Debug().
			Msg("Returning default format")
		return
	}

	log.Debug().
		Str("accept", rawHeader).
		Msg("Returning requested format")

	types := strings.Split(rawHeader, ",")
	for _, t := range types {
		semiPos := strings.Index(t, ";")
		if semiPos >= 0 {
			t = t[0:semiPos]
		}
		cleanType := strings.TrimSpace(t)
		if _, ok := supportedContentTypes[cleanType]; ok {
			contentType = cleanType
			return
		}
	}
	return
}
