package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/database/model"
	"golang.org/x/crypto/bcrypt"
)

func UpdateKeyboardStats(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")

	var user model.User
	var keyboardStats model.Keyboard

	// Fetch user from DB
	err := db.Find(&user, "id = ?", id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while fetching user from db",
			"data":    err,
		})
	}

	// Check if user was found
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to find the user",
			"data":    nil,
		})
	}

	if user.Username != "quantinium" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "WHY ARE YOU HERE???. WHO ARE YOU???. GO AWAY!!! SHOO SHOO",
		})
	}

	type authenticatedUpdateKeyboardStats struct {
		Password string `json:"password"`
		Keypress uint64 `json:"keypress"`
	}

	var UpdateKeyboardStats authenticatedUpdateKeyboardStats
	err = c.BodyParser(&UpdateKeyboardStats)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request. Body is wrong",
			"data":    err,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UpdateKeyboardStats.Password)); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to authenticate the request",
			"data":    err,
		})
	}

	// Fetch the existing keyboard stats by username
	err = db.First(&keyboardStats, "username = ?", user.Username).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch keyboard stats from the db",
			"data":    err,
		})
	}

	// Update the keypress count
	keyboardStats.Keypress += UpdateKeyboardStats.Keypress

	// Save the updated keyboard stats
	if err := db.Save(&keyboardStats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to update the keyboard stats",
			"data":    err,
		})
	}

	// Return success message
	return c.Status(200).JSON(fiber.Map{
		"status":  "success", // Change to success
		"message": "Successfully updated the keyboard stats",
		"data":    keyboardStats,
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

	if keyboardStats.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "keyboard stats not found",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success", // Change "error" to "success"
		"message": "Successfully fetched keyboard stats",
		"data":    keyboardStats, // You might want to return the data as well
	})
}

func CreateKeyboardStats(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")
	var user model.User

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

	if user.Username != "quantinium" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "WHY ARE YOU HERE???. WHO ARE YOU???. GO AWAY!!! SHOO SHOO",
		})
	}

	type authKeyboardReq struct {
		Password string `json:"password"`
		Keypress uint64 `json:"keypress"`
	}

	var authReq authKeyboardReq
	var keyboardStatsData model.Keyboard
	if err := c.BodyParser(&authReq); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request. Body is wrong",
			"data":    err,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authReq.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized request",
			"data":    err,
		})
	}

	keyboardStatsData.Keypress = authReq.Keypress
	keyboardStatsData.Username = user.Username
	if err := db.Create(&keyboardStatsData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add the data to the database",
			"data":    err,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "successfully added the keyboard stat",
		"data":    keyboardStatsData,
	})
}
