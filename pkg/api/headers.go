package api

import (
	"net/http"
)

// Headers godoc
// @Summary Headers
// @Description returns a JSON array with the request HTTP headers
// @Tags HTTP API
// @Accept json
// @Produce json
// @Router /headers [get]
// @Success 200 {object} api.ArrayResponse
func (a *Api) EchoHeadersHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "EchoHeadersHandler")
	defer span.End()
	a.JSONResponse(w, r, r.Header)
}
