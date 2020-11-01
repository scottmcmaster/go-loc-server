package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/scottmcmaster/go-loc-server/locserver/loader"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// StringHandler handles a request for an individual string.
type StringHandler struct {
	ST *loader.StringTable
}

func (h StringHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lang, accept, param := ExtractLang(req)
	vars := mux.Vars(req)

	tag, _ := language.MatchStrings(*h.ST.Matcher, lang, accept, param)
	p := message.NewPrinter(tag)

	str := p.Sprintf(vars["str"])
	log.Debug().
		Str("str", str).
		Str("cookie", lang).
		Str("accept", accept).
		Str("language_tag", tag.String()).
		Msg("Returning string")
	p.Printf(str)

	data := []byte(str)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
