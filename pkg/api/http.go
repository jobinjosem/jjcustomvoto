package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"

	"go.uber.org/zap"
)



func (a *Api) JSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
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

func WriteError(err error, w http.ResponseWriter, r *http.Request, status int, debug bool) {
	logMessage := fmt.Sprintf("Error serving request [%v]: %v", r, err)

	if debug {
		logMessage += fmt.Sprintf("\nRequest Headers: %+v", r.Header)
		logMessage += fmt.Sprintf("\nRequest Body: %+v", r.Body)
	}

	log.Printf(logMessage)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	errorMessage := make(map[string]interface{})
	errorMessage["error"] = fmt.Sprintf("%v", err)

	if debug {
		errorMessage["method"] = r.Method
		errorMessage["url"] = r.URL.String()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = []byte{}
		}
		errorMessage["request_body"] = string(body)
	}

	json.NewEncoder(w).Encode(errorMessage)
}