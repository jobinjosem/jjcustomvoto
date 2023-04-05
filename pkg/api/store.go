package api

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Store godoc
// @Summary Upload file
// @Description writes the posted content to disk at /data/hash and returns the SHA1 hash of the content
// @Tags HTTP API
// @Accept json
// @Produce json
// @Router /store [post]
// @Success 200 {object} api.MapResponse
func (a *Api) StoreWriteHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "storeWriteHandler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.ErrorResponse(w, r, span, "reading the request body failed", http.StatusBadRequest)
		return
	}

	hash := hash(string(body))
	err = os.WriteFile(path.Join(a.Config.DataPath, hash), body, 0644)
	if err != nil {
		a.Logger.Warn("writing file failed", zap.Error(err), zap.String("file", path.Join(a.Config.DataPath, hash)))
		a.ErrorResponse(w, r, span, "writing file failed", http.StatusInternalServerError)
		return
	}
	a.JSONResponseCode(w, r, map[string]string{"hash": hash}, http.StatusAccepted)
}

// Store godoc
// @Summary Download file
// @Description returns the content of the file /data/hash if exists
// @Tags HTTP API
// @Accept json
// @Produce plain
// @Param hash path string true "hash value"
// @Router /store/{hash} [get]
// @Success 200 {string} string "file"
func (a *Api) StoreReadHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "storeReadHandler")
	defer span.End()

	hash := mux.Vars(r)["hash"]
	content, err := os.ReadFile(path.Join(a.Config.DataPath, hash))
	if err != nil {
		a.Logger.Warn("reading file failed", zap.Error(err), zap.String("file", path.Join(a.Config.DataPath, hash)))
		a.ErrorResponse(w, r, span, "reading file failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(content))
}

func hash(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
