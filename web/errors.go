package web

import (
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)

	message := "internal server error"
	errorResponse(w, r, 500, message)
}

func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(w, r, 422, errors)
}
