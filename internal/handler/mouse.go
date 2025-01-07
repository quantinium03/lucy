package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/database/model"
	"golang.org/x/crypto/bcrypt"
)

func CreateMouseStats(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")
	var user model.User

	if err := db.Find(&user, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to query the database",
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

	type AuthMouseReq struct {
		Password    string `json:"password"`
		RightClick  uint64 `json:"rightClick"`
		LeftClick   uint64 `json:"leftClick"`
		MouseTravel uint64 `json:"mouseTravel"`
	}

	var authReq AuthMouseReq
	if err := c.BodyParser(&authReq); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request. Failed to parse the body",
			"data":    err,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),
		[]byte(authReq.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized Request",
			"data":    err,
		})
	}

	mouseStats := model.Mouse{
		Username:    user.Username,
		LeftClick:   authReq.LeftClick,
		RightClick:  authReq.RightClick,
		MouseTravel: authReq.MouseTravel,
	}

	if err := db.Create(&mouseStats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create mouse stats",
			"data":    err,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "successfully created mouse stats",
		"data":    mouseStats,
	})
}

func GetMouseStats(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")

	var user model.User
	if err := db.Find(&user, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find user from db",
			"data":    err,
		})
	}

	var mouseStats model.Mouse
	if err := db.First(&mouseStats, "username = ?", user.Username).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch mouse stats",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "successfully fetched mouse stats",
		"data":    mouseStats,
	})
}

func UpdateMouseStats(c *fiber.Ctx) error {
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
			"message": "User not found",
			"data":    nil,
		})
	}

	if user.Username != "quantinium" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "WHY ARE YOU HERE???. WHO ARE YOU???. GO AWAY!!! SHOO SHOO",
		})
	}

	type AuthMouseReq struct {
		Password    string `json:"password"`
		LeftClick   uint64 `json:"leftClick"`
		RightClick  uint64 `json:"rightClick"`
		MouseTravel uint64 `json:"mouseTravel"`
	}

	var authReq AuthMouseReq
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
			"message": "Unauthorized Request",
			"data":    err,
		})
	}

	var mouseStats model.Mouse
	if err := db.First(&mouseStats, "username = ?", user.Username).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch mouse stats from the db",
			"data":    err,
		})
	}

	mouseStats.LeftClick += authReq.LeftClick
	mouseStats.RightClick += authReq.RightClick
	mouseStats.MouseTravel += authReq.MouseTravel

	if err := db.Save(&mouseStats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update mouse stats in the db",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully updated mouse stats",
		"data":    mouseStats,
	})
}
