package api

import (
	"net/http"
	"os"
)

// Panic godoc
// @Summary Panic
// @Description crashes the process with exit code 255
// @Tags HTTP API
// @Router /panic [get]
func (a *Api) PanicHandler(w http.ResponseWriter, r *http.Request) {
	// a.Logger.Info("Panic command received")
	os.Exit(255)
}
