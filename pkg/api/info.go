package api

import (
	"net/http"

	"runtime"
	"strconv"

	"github.com/stefanprodan/podinfo/pkg/version"
)

// Info godoc
// @Summary Runtime information
// @Description returns the runtime information
// @Tags HTTP API
// @Accept json
// @Produce json
// @Success 200 {object} api.RuntimeResponse
// @Router /api/info [get]
func (a *Api) InfoHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "InfoHandler")
	defer span.End()

	data := RuntimeResponse{
		Hostname:     a.Config.Hostname,
		Version:      version.VERSION,
		Revision:     version.REVISION,
		Logo:         a.Config.UILogo,
		Color:        a.Config.UIColor,
		Message:      a.Config.UIMessage,
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		Runtime:      runtime.Version(),
		NumGoroutine: strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
		NumCPU:       strconv.FormatInt(int64(runtime.NumCPU()), 10),
	}

	a.JSONResponse(w, r, data)
}

type RuntimeResponse struct {
	Hostname     string `json:"hostname"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
	Message      string `json:"message"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	Runtime      string `json:"runtime"`
	NumGoroutine string `json:"num_goroutine"`
	NumCPU       string `json:"num_cpu"`
}
