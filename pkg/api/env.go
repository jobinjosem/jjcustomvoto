package api

import (
	"net/http"

	"os"
)

// Env godoc
// @Summary Environment
// @Description returns the environment variables as a JSON array
// @Tags HTTP API
// @Accept json
// @Produce json
// @Router /env [get]
// @Success 200 {object} api.ArrayResponse
func (a *Api) EnvHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "EnvHandler")
	defer span.End()
	a.JSONResponse(w, r, os.Environ())
}
