package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (app *application) readID(r *http.Request, key string) (int64, error) {
	postID := chi.URLParam(r, key)
	id, err := strconv.ParseInt(postID, 10, 64)
	return id, err
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *application) writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return app.writeJSON(w, status, &envelope{
		Error: message,
	})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	defer r.Body.Close()

	if err := dec.Decode(&data); err != nil {
		return err
	}
	return nil
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, v interface{}) error {
	type envelope struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}

	return app.writeJSON(w, status, envelope{
		Data: v,
	})
}

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
