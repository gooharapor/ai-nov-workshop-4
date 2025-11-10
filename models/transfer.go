package models

import (
	"time"

	"gorm.io/gorm"
)

// TransferStatus represents the status of a transfer
type TransferStatus string

const (
	TransferStatusPending    TransferStatus = "pending"
	TransferStatusProcessing TransferStatus = "processing"
	TransferStatusCompleted  TransferStatus = "completed"
	TransferStatusFailed     TransferStatus = "failed"
	TransferStatusCancelled  TransferStatus = "cancelled"
	TransferStatusReversed   TransferStatus = "reversed"
)

// Transfer represents a point transfer between users
type Transfer struct {
	ID             uint            `gorm:"primaryKey" json:"transferId,omitempty"`
	FromUserID     uint            `gorm:"not null;index:idx_transfers_from" json:"fromUserId"`
	ToUserID       uint            `gorm:"not null;index:idx_transfers_to" json:"toUserId"`
	Amount         int             `gorm:"not null;check:amount > 0" json:"amount"`
	Status         TransferStatus  `gorm:"not null;type:text" json:"status"`
	Note           string          `gorm:"type:text" json:"note,omitempty"`
	IdempotencyKey string          `gorm:"uniqueIndex;not null;size:128" json:"idemKey"`
	CreatedAt      time.Time       `gorm:"index:idx_transfers_created" json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	CompletedAt    *time.Time      `json:"completedAt,omitempty"`
	FailReason     string          `gorm:"type:text" json:"failReason,omitempty"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TransferCreateRequest for creating a new transfer
type TransferCreateRequest struct {
	FromUserID uint   `json:"fromUserId" binding:"required,min=1"`
	ToUserID   uint   `json:"toUserId" binding:"required,min=1"`
	Amount     int    `json:"amount" binding:"required,min=1"`
	Note       string `json:"note" binding:"max=512"`
}

// TransferResponse wraps transfer data
type TransferResponse struct {
	Transfer *Transfer `json:"transfer"`
}

// TransferListResponse for paginated transfer list
type TransferListResponse struct {
	Data     []Transfer `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int64      `json:"total"`
}
