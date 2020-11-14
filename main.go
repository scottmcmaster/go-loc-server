package main

import (
	"errors"
	"flag"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/scottmcmaster/go-loc-server/locserver/handlers"
	"github.com/scottmcmaster/go-loc-server/locserver/loader"
)

type loaderType string

const (
	goText loaderType = "gotext"
	xliff2            = "xliff2"
)

func (lt loaderType) IsValid() error {
	switch lt {
	case goText, xliff2:
		return nil
	}
	return errors.New("invalid loader type")
}

var lang = flag.String("lang", "en-us", "use language")
var localesDir = flag.String("localesdir", "./locales-gotext", "base directory of locale files")
var loaderTypeFl = flag.String("loader", "gotext", "loader type")
var debug = flag.Bool("debug", false, "sets log level to debug")
var server = flag.Bool("server", false, "starts in server mode")
var port = flag.Int("port", 3001, "http port")

func main() {
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	lt := loaderType(*loaderTypeFl)
	if err := lt.IsValid(); err != nil {
		log.Panic().Err(err).Msg("Loader type")
	}
	ldr, err := createLoader(lt)
	if err != nil {
		log.Panic().Err(err).Msg("Loader")
	}

	strs := &loader.StringTable{
		LocalesDir: *localesDir,
		Loader:     ldr,
	}

	err = strs.Load()
	if err != nil {
		log.Panic().Err(err).Msg("Fatal error while loading")
	}

	if *server {
		startServer(strs, *port)
	} else {
		log.Info().
			Str("language", *lang).
			Msg("Requesting language")

		tag, _ := language.MatchStrings(*strs.Matcher, *lang)
		log.Info().Str("language_tag", tag.String()).Msg("Using language tag")

		p := message.NewPrinter(tag)
		p.Printf("Hello world!")
	}
}

func createLoader(lt loaderType) (loader.Loader, error) {
	log.Info().Str("loaderType", string(lt)).Msg("Creating loader")

	switch lt {
	case goText:
		return loader.NewGoTextJSONLoader(), nil
	case xliff2:
		return loader.NewXLIFF2Loader(), nil
	}

	return nil, errors.New("unknown loader type " + string(lt))
}

func startServer(strs *loader.StringTable, port int) {
	mux := mux.NewRouter()

	sHandler := handlers.StringHandler{
		ST: strs,
	}
	mux.Handle("/v1/strings/{str}", sHandler)

	ssHandler := handlers.StringsHandler{
		ST: strs,
	}
	mux.Handle("/v1/strings", ssHandler)

	//Create the server.
	log.Info().
		Int("port", port).
		Msg("Starting server")

	s := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal().Int("port", port).Err(err).Msg("Failed to start server")
	}
}
