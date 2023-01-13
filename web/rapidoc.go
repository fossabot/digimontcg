package web

import (
	_ "embed"
	"net/http"
)

//go:embed template/rapidoc.html
var rapidocHTML []byte

func (app *Application) handleRapidoc(w http.ResponseWriter, r *http.Request) {
	w.Write(rapidocHTML)
}
