package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/quantinium03/lucy/internal/config"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/route"
)

func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	route.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./public/404.html")
	})


	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
