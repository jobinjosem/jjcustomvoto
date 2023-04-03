package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"github.com/jobinjosem/jjcustomvoto/emojivoto-web/web"
)



func (a *Api) JSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	if err != nil {
		web.WriteError(err, w, r, http.StatusInternalServerError, true)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON(body))
}

func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return out.Bytes()
}