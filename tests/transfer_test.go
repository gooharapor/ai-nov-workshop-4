package tests

import (
	"testing"

	"class-go-ai/models"
	"class-go-ai/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&models.User{}, &models.Transfer{}, &models.PointLedger{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateTransfer_Success(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewTransferService(db)

	// Create test users
	user1 := &models.User{Name: "Alice", Email: "alice@test.com", Points: 1000}
	user2 := &models.User{Name: "Bob", Email: "bob@test.com", Points: 500}
	db.Create(user1)
	db.Create(user2)

	// Create transfer
	req := &models.TransferCreateRequest{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Amount:     250,
		Note:       "Test transfer",
	}

	transfer, err := service.CreateTransfer(req)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if transfer.Status != models.TransferStatusCompleted {
		t.Errorf("Expected status completed, got: %s", transfer.Status)
	}

	if transfer.Amount != 250 {
		t.Errorf("Expected amount 250, got: %d", transfer.Amount)
	}

	// Verify points
	var updatedUser1, updatedUser2 models.User
	db.First(&updatedUser1, user1.ID)
	db.First(&updatedUser2, user2.ID)

	if updatedUser1.Points != 750 {
		t.Errorf("Expected user1 points 750, got: %d", updatedUser1.Points)
	}

	if updatedUser2.Points != 750 {
		t.Errorf("Expected user2 points 750, got: %d", updatedUser2.Points)
	}

	// Verify ledger entries
	var ledgers []models.PointLedger
	db.Find(&ledgers)

	if len(ledgers) != 2 {
		t.Errorf("Expected 2 ledger entries, got: %d", len(ledgers))
	}
}

func TestCreateTransfer_InsufficientPoints(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewTransferService(db)

	// Create test users
	user1 := &models.User{Name: "Alice2", Email: "alice2@test.com", Points: 100}
	user2 := &models.User{Name: "Bob2", Email: "bob2@test.com", Points: 500}
	if err := db.Create(user1).Error; err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	if err := db.Create(user2).Error; err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	t.Logf("User1 ID: %d, User2 ID: %d", user1.ID, user2.ID)

	// Try to create transfer with insufficient points
	req := &models.TransferCreateRequest{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Amount:     250,
	}

	transfer, err := service.CreateTransfer(req)

	// Should return error
	if err != services.ErrInsufficientPoints {
		t.Errorf("Expected ErrInsufficientPoints, got: %v", err)
	}

	// Transfer should be created but failed
	if transfer == nil {
		t.Error("Expected transfer object even on failure")
		return
	}

	if transfer.Status != models.TransferStatusFailed {
		t.Errorf("Expected status failed, got: %s", transfer.Status)
	}

	// Points should not change
	var updatedUser1 models.User
	db.First(&updatedUser1, user1.ID)

	if updatedUser1.Points != 100 {
		t.Errorf("Expected user1 points unchanged at 100, got: %d", updatedUser1.Points)
	}
}

func TestCreateTransfer_SameUser(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewTransferService(db)

	// Create test user
	user1 := &models.User{Name: "Alice3", Email: "alice3@test.com", Points: 1000}
	db.Create(user1)

	// Try to create transfer to same user
	req := &models.TransferCreateRequest{
		FromUserID: user1.ID,
		ToUserID:   user1.ID,
		Amount:     250,
	}

	_, err := service.CreateTransfer(req)

	// Should return error
	if err != services.ErrSameUser {
		t.Errorf("Expected ErrSameUser, got: %v", err)
	}
}

func TestGetTransferByIdemKey(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewTransferService(db)

	// Create test users
	user1 := &models.User{Name: "Alice4", Email: "alice4@test.com", Points: 1000}
	user2 := &models.User{Name: "Bob4", Email: "bob4@test.com", Points: 500}
	db.Create(user1)
	db.Create(user2)

	// Create transfer
	req := &models.TransferCreateRequest{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Amount:     250,
	}

	transfer, err := service.CreateTransfer(req)
	if err != nil {
		t.Fatalf("Failed to create transfer: %v", err)
	}

	if transfer == nil {
		t.Fatal("Transfer is nil")
	}

	// Get transfer by idemKey
	found, err := service.GetTransferByIdemKey(transfer.IdempotencyKey)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if found.ID != transfer.ID {
		t.Errorf("Expected transfer ID %d, got: %d", transfer.ID, found.ID)
	}
}

func TestGetTransfersByUserID(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewTransferService(db)

	// Create test users
	user1 := &models.User{Name: "Alice5", Email: "alice5@test.com", Points: 1000}
	user2 := &models.User{Name: "Bob5", Email: "bob5@test.com", Points: 1000}
	user3 := &models.User{Name: "Charlie5", Email: "charlie5@test.com", Points: 1000}
	db.Create(user1)
	db.Create(user2)
	db.Create(user3)

	// Create transfers
	service.CreateTransfer(&models.TransferCreateRequest{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Amount:     100,
	})

	service.CreateTransfer(&models.TransferCreateRequest{
		FromUserID: user3.ID,
		ToUserID:   user1.ID,
		Amount:     50,
	})

	// Get transfers for user1
	result, err := service.GetTransfersByUserID(user1.ID, 1, 20)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected 2 transfers for user1, got: %d", result.Total)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 transfer records, got: %d", len(result.Data))
	}
}
