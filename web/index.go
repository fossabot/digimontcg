package web

import (
	_ "embed"
	"net/http"
)

//go:embed template/index.html
var indexHTML []byte

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write(indexHTML)
}
