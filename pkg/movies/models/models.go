package models

import (
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
		go Db.AutoMigrate()
	}
}

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Overview    string   `json:"overview"`
	Genres      []string `json:"genres"`
	Poster      string   `json:"poster"`
	ReleaseDate int64    `json:"release_date"`
}

var MoviesIndex meilisearch.IndexManager
