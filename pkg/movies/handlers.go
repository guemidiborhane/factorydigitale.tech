package movies

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	_ "github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/movies/models"
	"github.com/meilisearch/meilisearch-go"
)

// @Summary	Index Movies
// @Tags		Core
// @Produce	json
// @Param		offset	query		int	false	"offset for paging"
// @Success	200		{object}	[]models.Movie
// @Failure	403		{object}	errors.HttpError
// @Failure	500		{object}	errors.HttpError
// @Router		/api/movies [get]
func IndexMovies(c *fiber.Ctx) error {
	var movies []models.Movie
	var results meilisearch.DocumentsResult
	offset := c.QueryInt("offset", 0)

	models.MoviesIndex.GetDocuments(&meilisearch.DocumentsQuery{
		Limit:  10,
		Offset: int64(offset),
	}, &results)

	jsonData, _ := json.Marshal(results.Results)
	json.Unmarshal(jsonData, &movies)

	return c.Status(fiber.StatusOK).JSON(movies)
}
