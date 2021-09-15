package main

import (
	"fiber_example/config"
	"fiber_example/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Ruta no definida 😉",
		})
	})

	api := app.Group("/v1")

	routes.SloganRoute(api.Group("/slogan"))
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error al cargar archivo .env")
		}
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	config.ConnectDB()

	setupRoutes(app)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error al iniciar la app")
		panic(err)
	}
}
