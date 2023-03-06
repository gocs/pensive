package router

import (
	"io"
	"log"
	"net/http"

	"github.com/gocs/pensive/internal/manager"
	"github.com/gorilla/mux"
)

func (a *App) GetObject(w http.ResponseWriter, r *http.Request) {
	self, err := manager.AuthSelf(r, a.session, a.client, UserIDSession)
	if err != nil {
		log.Println("unauthorized:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := self.Username(r.Context())
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	filename := mux.Vars(r)["filename"]
	if filename == "" {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}

	f, err := a.objs.GetObject(r.Context(), u, filename)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, f)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
}
