package handlers

import (
	"errors"
	"strconv"

	"class-go-ai/database"
	"class-go-ai/models"
	"class-go-ai/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var transferService *services.TransferService

// InitTransferService initializes the transfer service
func InitTransferService() {
	transferService = services.NewTransferService(database.DB)
}

// CreateTransfer handles POST /transfers
func CreateTransfer(c *fiber.Ctx) error {
	if transferService == nil {
		InitTransferService()
	}

	req := new(models.TransferCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Invalid input format",
		})
	}

	// Validate required fields
	if req.FromUserID == 0 || req.ToUserID == 0 || req.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "fromUserId, toUserId, and amount are required and must be greater than 0",
		})
	}

	transfer, err := transferService.CreateTransfer(req)
	
	if err != nil {
		switch {
		case errors.Is(err, services.ErrSameUser):
			return c.Status(422).JSON(fiber.Map{
				"error":   "INVALID_OPERATION",
				"message": "Cannot transfer to the same user",
			})
		case errors.Is(err, services.ErrUserNotFound):
			return c.Status(404).JSON(fiber.Map{
				"error":   "USER_NOT_FOUND",
				"message": "One or both users not found",
			})
		case errors.Is(err, services.ErrInsufficientPoints):
			return c.Status(409).JSON(fiber.Map{
				"error":   "INSUFFICIENT_POINTS",
				"message": "Sender does not have enough points",
			})
		default:
			return c.Status(500).JSON(fiber.Map{
				"error":   "INTERNAL_ERROR",
				"message": "Failed to create transfer",
			})
		}
	}

	// Set Idempotency-Key header
	c.Set("Idempotency-Key", transfer.IdempotencyKey)

	return c.Status(201).JSON(models.TransferResponse{
		Transfer: transfer,
	})
}

// GetTransfer handles GET /transfers/{id}
func GetTransfer(c *fiber.Ctx) error {
	if transferService == nil {
		InitTransferService()
	}

	idemKey := c.Params("id")
	if idemKey == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Transfer ID is required",
		})
	}

	transfer, err := transferService.GetTransferByIdemKey(idemKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error":   "NOT_FOUND",
				"message": "Transfer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfer",
		})
	}

	return c.JSON(models.TransferResponse{
		Transfer: transfer,
	})
}

// ListTransfers handles GET /transfers?userId=X&page=1&pageSize=20
func ListTransfers(c *fiber.Ctx) error {
	if transferService == nil {
		InitTransferService()
	}

	// Get userId from query (required)
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId query parameter is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil || userID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId must be a valid positive integer",
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	result, err := transferService.GetTransfersByUserID(uint(userID), page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfers",
		})
	}

	return c.JSON(result)
}
