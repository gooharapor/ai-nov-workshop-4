package handlers

import (
	"class-go-ai/database"
	"class-go-ai/models"

	"github.com/gofiber/fiber/v2"
)

// GetUsers returns all users
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	
	result := database.DB.Order("id DESC").Find(&users)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

// GetUser returns a single user by ID
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
	input := new(models.UserInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Validate required fields
	if input.Name == "" || input.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name and email are required",
		})
	}

	user := models.User{
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Address: input.Address,
		Avatar:  input.Avatar,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(user)
}

// UpdateUser updates an existing user
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Check if user exists
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	input := new(models.UserInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Update user fields
	updates := map[string]interface{}{
		"name":    input.Name,
		"email":   input.Email,
		"phone":   input.Phone,
		"address": input.Address,
		"avatar":  input.Avatar,
	}

	result = database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Fetch updated user
	database.DB.First(&user, id)

	return c.JSON(user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Check if user exists
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Delete user
	result = database.DB.Delete(&user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
