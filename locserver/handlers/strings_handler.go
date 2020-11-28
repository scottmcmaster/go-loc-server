package handlers

import (
	"encoding/csv"
	"encoding/json"
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

type stringTranslation struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
}

func (h StringsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lang, acceptLang, param := ExtractLang(req)
	contentType := ExtractContentType(req)
	keyFilter := GetQueryParam(req, "kf")

	tag, _ := language.MatchStrings(*h.ST.Matcher, param, lang, acceptLang)

	log.Debug().
		Str("cookie", lang).
		Str("accept", acceptLang).
		Str("content_type", contentType).
		Str("key_filter", keyFilter).
		Str("language_tag", tag.String()).
		Msg("Returning strings")

	strs, err := h.ST.Loader.StringsByTag(tag)
	if err != nil {
		log.Error().Str("language_tag", tag.String()).Err(err).
			Msg("Getting strings for tag")
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("404 - Language not found"))
		return
	}

	switch contentType {
	case "application/json":
		err = writeJSON(res, strs, keyFilter)
	case "text/csv":
		fallthrough
	default:
		err = writeCSV(res, strs, keyFilter)
	}

	if err != nil {
		log.Error().Str("language_tag", tag.String()).Err(err).
			Msg("Unexpected error writing response")
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("500 - Unexpected error"))
		return
	}
}

func writeCSV(res http.ResponseWriter, strs *loader.StringCatalog, keyFilter string) error {
	res.Header().Set("Content-Type", "text/csv")
	w := csv.NewWriter(res)
	defer w.Flush()
	for k, v := range strs.Strings {
		if strings.HasPrefix(k, keyFilter) {
			w.Write([]string{k, v})
		}
	}
	return nil
}

func writeJSON(res http.ResponseWriter, strs *loader.StringCatalog, keyFilter string) error {
	res.Header().Set("Content-Type", "application/json")
	data := []stringTranslation{}
	for k, v := range strs.Strings {
		if strings.HasPrefix(k, keyFilter) {
			data = append(data, stringTranslation{
				Translation: v,
				ID:          k,
			})
		}
	}

	return json.NewEncoder(res).Encode(data)
}
