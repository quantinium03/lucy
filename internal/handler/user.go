package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/database/model"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	db := database.DB.DB
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad request. Body is wrong",
			"data":    err,
		})
	}

	// hash the password for storing the database
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user. failed to generate password hash",
			"data":    err,
		})
	}
	user.Password = string(hash[:])

	// create the user
	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user created successfully",
		"data":    user,
	})
}

func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")

	var user model.User

	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "user not found",
			"data":    nil,
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"status":  "success",
		"message": "user found",
		"data":    user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB.DB
	id := c.Params("id")

	var user model.User

	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "user not found",
			"data":    nil,
		})
	}

	res := db.Where("id = ?", id).Delete(&user)
	if res.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "failed to delete user",
			"data" : res.Error,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted successfully",
	})
}
