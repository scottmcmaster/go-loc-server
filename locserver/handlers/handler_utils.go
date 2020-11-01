package handlers

import "net/http"

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
