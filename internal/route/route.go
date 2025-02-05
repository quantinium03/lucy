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

	v1.Post("/keyboard/:id", handler.CreateKeyboardStats)
	v1.Get("/keyboard/:id", handler.GetKeyboardStats)
	v1.Put("/keyboard/:id", handler.UpdateKeyboardStats)

	v1.Post("/mouse/:id", handler.CreateMouseStats)
	v1.Get("/mouse/:id", handler.GetMouseStats)
	v1.Put("/mouse/:id", handler.UpdateMouseStats)

	v1.Get("/spotify/:id", handler.GetCurrentlyPlaying)
}
