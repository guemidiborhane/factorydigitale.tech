package models

import (
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
