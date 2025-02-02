package models

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Setup(database *gorm.DB) {
	Db = database

	MoviesIndex = storage.MeiliClient.Index("movies")
	filterableAttributes := []string{"genres", "id", "title", "overview"}
	MoviesIndex.UpdateFilterableAttributes(&filterableAttributes)
	if !fiber.IsChild() {
		go Db.AutoMigrate(&Favourite{})
	}
}

type Movie struct {
	ID           int      `json:"id" binding:"required"`
	Title        string   `json:"title" binding:"required"`
	Overview     string   `json:"overview" binding:"required"`
	Genres       []string `json:"genres" binding:"required"`
	Poster       string   `json:"poster" binding:"required"`
	ReleaseDate  int64    `json:"release_date" binding:"required"`
	InFavourites bool     `json:"in_favourites" binding:"required"`
}

type Favourite struct {
	ID      int `json:"id" binding:"required"`
	MovieID int `json:"movie_id" gorm:"uniqueIndex:idx_user_movie" binding:"required"`
	UserID  int `json:"user_id" gorm:"uniqueIndex:idx_user_movie" binding:"required"`
}

var MoviesIndex meilisearch.IndexManager

func (f *Favourite) Movie() (Movie, error) {
	var movie Movie
	if err := MoviesIndex.GetDocument(strconv.Itoa(f.MovieID), &meilisearch.DocumentQuery{
		Fields: []string{"id"},
	}, &movie); err != nil {
		return Movie{}, err
	}
	return movie, nil
}

func GetMoviesByIDs(ids []int) ([]Movie, error) {
	// Create a filter query to match documents with the given IDs
	filter := fmt.Sprintf("id IN [%s]", formatIDsForFilter(ids))

	// Set up the search parameters
	searchRequest := &meilisearch.SearchRequest{
		Filter: []string{filter},
	}

	// Perform the search
	results, err := MoviesIndex.Search("", searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	// Parse the results into Movie structs
	var movies []Movie

	jsonData, _ := json.Marshal(results.Hits)
	json.Unmarshal(jsonData, &movies)

	return movies, nil
}

// formatIDsForFilter formats the IDs array for the Meilisearch filter syntax
func formatIDsForFilter(ids []int) string {
	if len(ids) == 0 {
		return ""
	}

	// Build the comma-separated string of quoted IDs
	var result string
	for i, id := range ids {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf(`"%d"`, id)
	}
	return result
}
