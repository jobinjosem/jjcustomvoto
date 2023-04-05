package api

import (
	"html/template"
	"net/http"
	"path"
)

// Index godoc
// @Summary Index
// @Description renders podinfo UI
// @Tags HTTP API
// @Produce html
// @Router / [get]
// @Success 200 {string} string "OK"
func (a *Api) IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, span := a.Tracer.Start(r.Context(), "IndexHandler")
	defer span.End()

	tmpl, err := template.New("vue.html").ParseFiles(path.Join(a.Config.UIPath, "vue.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(path.Join(a.Config.UIPath, "vue.html") + err.Error()))
		return
	}

	data := struct {
		Title string
		Logo  string
	}{
		Title: a.Config.Hostname,
		Logo:  a.Config.UILogo,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, path.Join(a.Config.UIPath, "vue.html")+err.Error(), http.StatusInternalServerError)
	}
}
