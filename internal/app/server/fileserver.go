package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewFileHandler(fileRouter *mux.Router) {
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	fileRouter.PathPrefix("/media/").Handler(fileServer)
}
