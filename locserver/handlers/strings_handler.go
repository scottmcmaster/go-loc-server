package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"

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

	tag, _ := language.MatchStrings(*h.ST.Matcher, lang, acceptLang, param)

	log.Debug().
		Str("cookie", lang).
		Str("accept", acceptLang).
		Str("content_type", contentType).
		Str("language_tag", tag.String()).
		Msg("Returning strings")

	strs, err := h.ST.Loader.StringsByTag(tag)
	if err != nil {
		log.Error().Str("language_tag", tag.String()).Err(err).
			Msg("Getting strings for tag")
		res.Write([]byte("404 - Language not found"))
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch contentType {
	case "application/json":
		err = writeJSON(res, strs)
	case "text/csv":
		fallthrough
	default:
		err = writeCSV(res, strs)
	}

	if err != nil {
		log.Error().Str("language_tag", tag.String()).Err(err).
			Msg("Unexpected error writing response")
		res.Write([]byte("500 - Unexpected error"))
		res.WriteHeader(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)
}

func writeCSV(res http.ResponseWriter, strs *loader.StringCatalog) error {
	res.Header().Set("Content-Type", "text/csv")
	w := csv.NewWriter(res)
	for k, v := range strs.Strings {
		w.Write([]string{k, v})
	}
	return nil
}

func writeJSON(res http.ResponseWriter, strs *loader.StringCatalog) error {
	res.Header().Set("Content-Type", "application/json")
	data := []stringTranslation{}
	for k, v := range strs.Strings {
		data = append(data, stringTranslation{
			Translation: v,
			ID:          k,
		})
	}

	return json.NewEncoder(res).Encode(data)
}
