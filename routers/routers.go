package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thn-lee/go-fiber-101/movie"
)

func SetupRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")

	v1Group := apiGroup.Group("/v1")
	{
		v1Group.Get("/movie", movie.GetMovies)
		v1Group.Get("/movie/:id", movie.GetMovieById)
		v1Group.Post("/movie", movie.CreateMovie)
		v1Group.Delete("/movie/:id", movie.DeleteMovie)
		v1Group.Patch("/movie/:id", movie.UpdateMovie)
	}

}
