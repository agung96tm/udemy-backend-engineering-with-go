package main

import "net/http"

// healthCheckHandler godoc
//
//	@Summary		Health check
//	@Description	Mengecek status kesehatan API. Mengembalikan status, env, dan version.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"status (ok), env, version"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	_ = app.writeJSON(w, http.StatusOK, data)
}
