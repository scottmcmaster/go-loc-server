package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/scottmcmaster/go-loc-server/locserver/loader"
	"golang.org/x/text/language"
)

// StringsHandler handles a request for a full string catalog by language.
type StringsHandler struct {
	ST *loader.StringTable
}

func (h StringsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lang, accept, param := ExtractLang(req)

	tag, _ := language.MatchStrings(*h.ST.Matcher, lang, accept, param)

	sb := strings.Builder{}
	log.Debug().
		Str("cookie", lang).
		Str("accept", accept).
		Str("language_tag", tag.String()).
		Msg("Returning strings")

	res.Header().Set("Content-Type", "text/plain")

	strs, err := h.ST.Loader.StringsByTag(tag)
	if err != nil {
		log.Error().Str("languagetag", tag.String()).
			Msg("Getting strings for tag")
		res.Write([]byte("404 - Language not found"))
		res.WriteHeader(http.StatusNotFound)
		return
	}

	for k, v := range strs.Strings {
		sb.WriteString(fmt.Sprintf("%s\t%s\n", k, v))
	}

	data := []byte(sb.String())
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
