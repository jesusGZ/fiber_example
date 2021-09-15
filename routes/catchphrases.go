package routes

import (
	"fiber_example/controllers"

	"github.com/gofiber/fiber/v2"
)

func SloganRoute(route fiber.Router) {
	route.Get("/", controllers.GetAllSlogans)
	route.Get("/:id", controllers.GetSlogan)
	route.Post("/", controllers.AddSlogan)
	route.Put("/:id", controllers.UpdateSlogan)
	route.Delete("/:id", controllers.DeleteSlogan)
}
