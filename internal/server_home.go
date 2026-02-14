package internal

import (
	"embed"
	"net/http"
)

//go:embed index.html
var template embed.FS

func (s *server) handlerHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, template, "index.html")
}
