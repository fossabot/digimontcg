package web

import (
	"log"
	"net/http"
)

func (app *Application) handleCards(w http.ResponseWriter, r *http.Request) {
	err := writeJSON(w, 200, envelope{"cards": app.cards})
	if err != nil {
		log.Println(err)

		code := 500
		http.Error(w, http.StatusText(code), code)
		return
	}
}
