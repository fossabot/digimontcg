package web

import (
	"log"
	"net/http"
)

func (app *Application) handleSets(w http.ResponseWriter, r *http.Request) {
	err := writeJSON(w, 200, envelope{"sets": app.sets})
	if err != nil {
		log.Println(err)

		code := 500
		http.Error(w, http.StatusText(code), code)
		return
	}
}
