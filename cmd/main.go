package main

import (
	"Redioteka/internal/pkg/domain"
	"encoding/json"
	"fmt"
)


func main() {
	//server.RunServer(":8081")
	res, _ := json.Marshal(domain.Movie{
		ID:          1,
		Title:       "Film",
		Description: "Test data",
		Rating:      9,
		Countries:   []string{"Japan", "South Korea"},
		IsFree:      false,
		Genres:      []string{"Comedy"},
		Actors:      []string{"Sana", "Momo", "Mina"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "2012",
		Director:    []string{"James Cameron"},
	})
	fmt.Println(string(res))
}
