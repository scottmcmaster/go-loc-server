package main

import (
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

var lang = flag.String("lang", "en-us", "use language")
var localesDir = flag.String("localesdir", "./locales", "base directory of locale files")
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

	strs := &loader.StringTable{
		LocalesDir: *localesDir,
		Loader:     loader.NewGoTextJSONLoader(),
	}

	err := strs.Load()
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
