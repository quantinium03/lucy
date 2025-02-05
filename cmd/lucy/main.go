package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/quantinium03/lucy/internal/config"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/route"
	"github.com/quantinium03/lucy/internal/util"
)

func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	go util.FetchSpotifyData()
	go util.GetAccessToken()

	route.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
