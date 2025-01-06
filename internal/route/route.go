package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantinium03/lucy/internal/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/user/:id", handler.GetSingleUser)
	v1.Post("/user", handler.CreateUser)
	v1.Delete("/user/:id", handler.DeleteUser)

	v1.Post("/keyboard/:id", handler.UpdateKeyboardStats)
	v1.Get("/keyboard/:id", handler.GetKeyboardStats)
}
