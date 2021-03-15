package movie

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Rating      int      `json:"rating"`
	Countries   []string `json:"countries"`
	IsFree      bool     `json:"is_free"`
	Genres      []string `json:"genres"`
	Actors      []string `json:"actors"`
	Avatar      string   `json:"movie_avatar,omitempty"`
	Type        string   `json:"type"`
	Year        string   `json:"year"`
	Director    []string `json:"director"`
}

type moviesData struct {
	sync.Mutex
	movies map[uint]*Movie
}

func (data *moviesData) fill() {
	data.movies[1] = &Movie{
		ID:          1,
		Title:       "Film",
		Description: "Test data",
		Rating:      9,
		Countries:   []string{"Japan", "South Korea"},
		IsFree:      false,
		Genres:      []string{"Comedy"},
		Actors:      []string{"Sana", "Momo", "Mina"},
		Avatar:      "/static/movies/default.jpg",
		Type:        "movie",
		Year:        "2012",
		Director:    []string{"James Cameron"},
	}
}

func (data *moviesData) getByID(id uint) *Movie {
	data.Lock()
	movie, exists := data.movies[id]
	data.Unlock()
	if exists {
		return movie
	}
	return nil
}

var data = moviesData{
	movies: make(map[uint]*Movie),
}

func init() {
	data.fill()
}

type Handler struct {
}

//Get - handler for viewing movie info page
func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Print("Id is not a number")
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	movie := data.getByID(uint(id))
	if movie == nil {
		log.Printf("This movie does not exist")
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
}
