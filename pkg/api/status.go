package api

import (
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

// Status godoc
// @Summary Status code
// @Description sets the response status code to the specified code
// @Tags HTTP API
// @Accept json
// @Produce json
// @Param code path int true "status code to return"
// @Router /status/{code} [get]
// @Success 200 {object} api.MapResponse
func (a *Api) StatusHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "StatusHandler")
	defer span.End()

	vars := mux.Vars(r)

	code, err := strconv.Atoi(vars["code"])
	if err != nil {
		a.ErrorResponse(w, r, span, err.Error(), http.StatusBadRequest)
		return
	}

	a.JSONResponseCode(w, r, map[string]int{"status": code}, code)
}
