package moviegen

import (
	"Redioteka/internal/pkg/domain"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"strconv"
	"time"
)

var countries = []string{
	"Dominica",
	"Barbados",
	"Puerto Rico",
	"Micronesia",
	"Ukraine",
	"Vanuatu",
	"Congo",
	"Uruguay",
	"French Polynesia",
	"China",
}

var actors = []string{
	"Twila Smitham",
	"Hunter Stehr",
	"Ansel Lubowitz",
	"Emilio Bogan",
	"Willie Christiansen",
	"Eve Rath",
	"Mathew Crooks",
	"Anna Hammes",
	"Alta Kuhic",
	"Cody Beer",
	"Emmitt Bruen",
	"Omer Herman",
	"Hertha Klocko",
	"Julio Quigley",
	"Stefanie Waters",
	"Dorris Beahan",
	"Stacy Schamberger",
	"Kiley OReilly",
	"Coty Erdman",
	"Mireya Batz",
}

var genres = []string{
	"horror",
	"sad",
	"happy",
	"nice",
	"thriller",
	"i don't know this genre",
	"another genre name",
	"iron man?",
}

var directors = []string{
	"James Cameron",
	"Stephen Spielberg",
	"Mikhalkov",
}

var types = []domain.MovieType{
	domain.MovieT,
	domain.SeriesT,
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func randomSlice(source []string, maxCount int) []string {
	rand.Shuffle(len(source), func(i, j int) {
		source[i], source[j] = source[j], source[i]
	})
	limit := Abs(rand.Int()) % (maxCount)
	if limit == 0 {
		limit = 1
	}
	return source[:limit]
}

func randomFillActors(m *domain.Movie) {
	m.Actors = randomSlice(actors, len(actors))
}

func randomFillCountries(m *domain.Movie) {
	m.Countries = randomSlice(countries, 2)
}

func randomFillGenres(m *domain.Movie) {
	m.Genres = randomSlice(genres, 3)
}

func randomFillDirectors(m *domain.Movie) {
	m.Director = randomSlice(directors, 2)
}

func RandomMovie(id uint) domain.Movie {
	rand.Seed(time.Now().UnixNano())
	faker := gofakeit.New(0)
	m := domain.Movie{}
	faker.Struct(&m)
	randomFillActors(&m)
	randomFillCountries(&m)
	randomFillGenres(&m)
	randomFillDirectors(&m)
	m.Year = strconv.Itoa(faker.Year())
	m.Type = types[faker.Number(0, 1)]
	m.ID = id
	fmt.Println("generated movie with id", id)
	return m
}
