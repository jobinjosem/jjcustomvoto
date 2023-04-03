package api

import (
	"net/http"

	"github.com/jobinjosem/jjcustomvoto/pkg/version"
)

// Version godoc
// @Summary Version
// @Description returns podinfo version and git commit hash
// @Tags HTTP API
// @Produce json
// @Router /version [get]
// @Success 200 {object} api.MapResponse
func (a *Api) VersionHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{
		"version": version.VERSION,
		"commit":  version.REVISION,
	}
	a.JSONResponse(w, r, result)
}
