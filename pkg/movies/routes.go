package movies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users"
)

func setupRoutes(r fiber.Router) {
	r.Get("/movies", users.Can("movies:index"), IndexMovies)
	r.Post("/movies/favourite", users.Can("movies:favourite"), FavouriteMovie)
	r.Get("/movies/favourites", users.Can("movies:favourite"), IndexFavourites)
}
