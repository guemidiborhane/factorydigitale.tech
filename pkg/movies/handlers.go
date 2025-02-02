package movies

import (
	"encoding/json"
	"slices"

	e "errors"

	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	_ "github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/movies/models"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

// @Summary	Index Movies
// @Tags		Movie
// @Produce	json
// @Param		offset	query		int	false	"offset for paging"
// @Success	200		{object}	[]models.Movie
// @Failure	403		{object}	errors.HttpError
// @Failure	500		{object}	errors.HttpError
// @Router		/api/movies [get]
func IndexMovies(c *fiber.Ctx) error {
	var movies []models.Movie
	var results meilisearch.DocumentsResult
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		return err
	}
	offset := c.QueryInt("offset", 0)

	models.MoviesIndex.GetDocuments(&meilisearch.DocumentsQuery{
		Limit:  10,
		Offset: int64(offset),
	}, &results)

	jsonData, _ := json.Marshal(results.Results)
	json.Unmarshal(jsonData, &movies)

	favourites, err := user.Favourites()
	if err != nil {
		return err
	}

	for i, movie := range movies {
		movies[i].InFavourites = slices.Contains(favourites, movie.ID)
	}

	return c.Status(fiber.StatusOK).JSON(movies)
}

type FavouriteRequestParams struct {
	MovieID int `json:"movie_id"`
}

// @Summary	Toggle movie to/from your Favourites
// @Tags		Movie
// @Produce	json
// @Param		body	body		FavouriteRequestParams	true "Favourite Request"
// @Success	200		{object}	models.Favourite
// @Failure	403		{object}	errors.HttpError
// @Failure	500		{object}	errors.HttpError
// @Router		/api/movies/favourite [post]
func FavouriteMovie(c *fiber.Ctx) error {
	var body FavouriteRequestParams
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&body); err != nil {
		return errors.BadRequest(err.Error())
	}
	favourite := &models.Favourite{
		UserID:  user.ID,
		MovieID: body.MovieID,
	}

	if err := models.Db.Create(&favourite).Error; err != nil {
		if e.Is(err, gorm.ErrDuplicatedKey) {
			if err := models.Db.First(&favourite, "movie_id = ? AND user_id = ?", favourite.MovieID, favourite.UserID).Error; err != nil {
				return errors.Unexpected(err.Error())
			}

			logger.Debug("find", logger.Attrs{
				"favourite": favourite,
			})

			if err := models.Db.Delete(&favourite).Error; err != nil {
				return errors.Unexpected(err.Error())
			}

			favourite = &models.Favourite{}
		}
	}

	return c.Status(fiber.StatusOK).JSON(favourite)
}
