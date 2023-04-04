package api

import "net/http"

func (a *Api) ConfigReadHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "configReadHandler")
	defer span.End()

	files := make(map[string]string)
	if watcher != nil {
		watcher.Cache.Range(func(key interface{}, value interface{}) bool {
			files[key.(string)] = value.(string)
			return true
		})
	}

	a.JSONResponse(w, r, files)
}
