package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewFileHandler(fileRouter *mux.Router) {
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./img")))
	fileRouter.PathPrefix("/movies/").Handler(fileServer)
	fileRouter.PathPrefix("/actors/").Handler(fileServer)
	fileRouter.PathPrefix("/users/").Handler(fileServer)
}
