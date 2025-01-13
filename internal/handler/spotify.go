package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/database/model"
)

func GetCurrentlyPlaying(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")

	var user model.User

	if err := db.Find(&user, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while fetching user from db",
			"data":    err,
		})
	}

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "user not found",
			"data":    nil,
		})
	}

	if user.Username != "quantinium" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "WHO ARE YOU BRO??? YOU ARE NOT ALLOWED HERE!! GO AWAAAAAY",
		})
	}

	var spotify model.Spotify
	if err := db.Find(&spotify, "username = ?", user.Username).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldnt fetch the url from the database",
			"data":    err,
		})
	}

	if spotify.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "spotify data not found",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully fetched the spotify data",
		"data":    spotify,
	})
}

