package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "modernc.org/sqlite"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserInput for create/update operations (without ID and timestamps)
type UserInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Avatar  string `json:"avatar"`
}

var db *sql.DB

// Initialize database and create users table
func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		return err
	}

	// Create users table
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		phone TEXT,
		address TEXT,
		avatar TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// Get all users
func getUsers(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, name, email, phone, address, avatar, created_at, updated_at FROM users ORDER BY id DESC")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to scan user",
			})
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

// Get user by ID
func getUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user User
	err := db.QueryRow("SELECT id, name, email, phone, address, avatar, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

// Create new user
func createUser(c *fiber.Ctx) error {
	input := new(UserInput)

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

	result, err := db.Exec(
		"INSERT INTO users (name, email, phone, address, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		input.Name, input.Email, input.Phone, input.Address, input.Avatar, time.Now(), time.Now(),
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	id, _ := result.LastInsertId()

	// Fetch the created user
	var user User
	db.QueryRow("SELECT id, name, email, phone, address, avatar, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	return c.Status(201).JSON(user)
}

// Update user
func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	input := new(UserInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Check if user exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&exists)
	if err != nil || exists == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update user
	_, err = db.Exec(
		"UPDATE users SET name = ?, email = ?, phone = ?, address = ?, avatar = ?, updated_at = ? WHERE id = ?",
		input.Name, input.Email, input.Phone, input.Address, input.Avatar, time.Now(), id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Fetch updated user
	var user User
	db.QueryRow("SELECT id, name, email, phone, address, avatar, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	return c.JSON(user)
}

// Delete user
func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if user exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&exists)
	if err != nil || exists == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Delete user
	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	app := fiber.New()

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hello world",
		})
	})

	// User routes
	app.Get("/users", getUsers)
	app.Get("/users/:id", getUserByID)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	log.Println("Server starting on port 3000...")
	app.Listen(":3000")
}
