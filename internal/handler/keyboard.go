package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/database/model"
)

func UpdateKeyboardStats(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")

	var user model.User

	err := db.Find(&user, "id = ?", id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while fetching user from db",
			"data":    err,
		})
	}

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to find the user",
			"data":    err,
		})
	}

	var keyboardStats model.Keyboard
	err = c.BodyParser(&keyboardStats)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request. Body is wrong",
			"data":    err,
		})
	}

	type UpdateKeyboardStats struct {
		Keypress uint64 `json:"keypress"`
	}

	var updateKeyboardStats UpdateKeyboardStats
	err = db.Find(&updateKeyboardStats, "username = ?", user.Username).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch the keyboard stats from the db",
			"data":    err,
		})
	}

	keyboardStats.Keypress += updateKeyboardStats.Keypress
	err = db.Save(&keyboardStats).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to update the keyboard stats",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "error",
		"message": "Successfully updated the keyboard stats",
	})

}

func GetKeyboardStats(c *fiber.Ctx) error {
	var user model.User
	id := c.Params("id")
	db := database.DB.DB

	if err := db.Find(&user, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch user from db",
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

	var keyboardStats model.Keyboard
	if err := db.Find(&keyboardStats, "username = ?", user.Username).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch keyboard stats from db",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "error",
		"message": "Successfully fetched keyboard stats",
	})
}
