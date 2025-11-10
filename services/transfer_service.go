package services

import (
	"errors"
	"fmt"
	"time"

	"class-go-ai/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInsufficientPoints = errors.New("insufficient points")
	ErrSameUser           = errors.New("cannot transfer to the same user")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAmount      = errors.New("amount must be greater than 0")
)

// TransferService handles business logic for transfers
type TransferService struct {
	db *gorm.DB
}

// NewTransferService creates a new transfer service
func NewTransferService(db *gorm.DB) *TransferService {
	return &TransferService{db: db}
}

// CreateTransfer creates a new transfer with atomic transaction
func (s *TransferService) CreateTransfer(req *models.TransferCreateRequest) (*models.Transfer, error) {
	// Validation
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if req.FromUserID == req.ToUserID {
		return nil, ErrSameUser
	}

	// Generate idempotency key
	idemKey := uuid.New().String()

	transfer := &models.Transfer{
		FromUserID:     req.FromUserID,
		ToUserID:       req.ToUserID,
		Amount:         req.Amount,
		Note:           req.Note,
		IdempotencyKey: idemKey,
		Status:         models.TransferStatusPending,
	}

	// Start transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get sender
		var fromUser models.User
		if err := tx.First(&fromUser, req.FromUserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		// Get receiver
		var toUser models.User
		if err := tx.First(&toUser, req.ToUserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		// Check sufficient points
		if fromUser.Points < req.Amount {
			transfer.Status = models.TransferStatusFailed
			transfer.FailReason = "Insufficient points"
			if err := tx.Create(transfer).Error; err != nil {
				return err
			}
			return ErrInsufficientPoints
		}

		// Update status to processing
		transfer.Status = models.TransferStatusProcessing
		if err := tx.Create(transfer).Error; err != nil {
			return err
		}

		// Deduct points from sender
		fromUser.Points -= req.Amount
		if err := tx.Save(&fromUser).Error; err != nil {
			return err
		}

		// Add points to receiver
		toUser.Points += req.Amount
		if err := tx.Save(&toUser).Error; err != nil {
			return err
		}

		// Create ledger entries
		now := time.Now()

		// Sender ledger
		senderLedger := &models.PointLedger{
			UserID:       fromUser.ID,
			Change:       -req.Amount,
			BalanceAfter: fromUser.Points,
			EventType:    models.EventTypeTransferOut,
			TransferID:   &transfer.ID,
			Reference:    fmt.Sprintf("Transfer to user %d", toUser.ID),
			CreatedAt:    now,
		}
		if err := tx.Create(senderLedger).Error; err != nil {
			return err
		}

		// Receiver ledger
		receiverLedger := &models.PointLedger{
			UserID:       toUser.ID,
			Change:       req.Amount,
			BalanceAfter: toUser.Points,
			EventType:    models.EventTypeTransferIn,
			TransferID:   &transfer.ID,
			Reference:    fmt.Sprintf("Transfer from user %d", fromUser.ID),
			CreatedAt:    now,
		}
		if err := tx.Create(receiverLedger).Error; err != nil {
			return err
		}

		// Mark transfer as completed
		completedAt := time.Now()
		transfer.Status = models.TransferStatusCompleted
		transfer.CompletedAt = &completedAt
		if err := tx.Save(transfer).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ErrInsufficientPoints) {
			return transfer, err
		}
		return nil, err
	}

	return transfer, nil
}

// GetTransferByIdemKey retrieves a transfer by idempotency key
func (s *TransferService) GetTransferByIdemKey(idemKey string) (*models.Transfer, error) {
	var transfer models.Transfer
	if err := s.db.Where("idempotency_key = ?", idemKey).First(&transfer).Error; err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetTransfersByUserID retrieves transfers for a user with pagination
func (s *TransferService) GetTransfersByUserID(userID uint, page, pageSize int) (*models.TransferListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}

	var transfers []models.Transfer
	var total int64

	// Count total
	s.db.Model(&models.Transfer{}).
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Count(&total)

	// Get paginated results
	offset := (page - 1) * pageSize
	err := s.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transfers).Error

	if err != nil {
		return nil, err
	}

	return &models.TransferListResponse{
		Data:     transfers,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
