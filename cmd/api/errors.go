package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	_ = app.writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	_ = app.writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	_ = app.writeJSONError(w, http.StatusNotFound, err.Error())
}
