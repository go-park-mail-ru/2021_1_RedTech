package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewFileHandler(fileRouter *mux.Router) {
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	fileRouter.PathPrefix("/media/").Handler(fileServer)
}
