package models

import (
	"time"
)

// EventType represents the type of point ledger event
type EventType string

const (
	EventTypeTransferOut EventType = "transfer_out"
	EventTypeTransferIn  EventType = "transfer_in"
	EventTypeAdjust      EventType = "adjust"
	EventTypeEarn        EventType = "earn"
	EventTypeRedeem      EventType = "redeem"
)

// PointLedger represents an entry in the point ledger (append-only)
type PointLedger struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index:idx_ledger_user" json:"userId"`
	Change       int       `gorm:"not null" json:"change"` // +receive / -send
	BalanceAfter int       `gorm:"not null" json:"balanceAfter"`
	EventType    EventType `gorm:"not null;type:text" json:"eventType"`
	TransferID   *uint     `gorm:"index:idx_ledger_transfer" json:"transferId,omitempty"`
	Reference    string    `gorm:"type:text" json:"reference,omitempty"`
	Metadata     string    `gorm:"type:text" json:"metadata,omitempty"` // JSON text
	CreatedAt    time.Time `gorm:"index:idx_ledger_created" json:"createdAt"`
}
