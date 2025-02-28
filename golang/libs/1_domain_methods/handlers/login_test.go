package handlers

import (
	"de-net/libs/1_domain_methods/helpers"
	"de-net/libs/2_generated_models/model"
	"de-net/libs/3_infrastructure/test_sqlite"
	"de-net/libs/4_common/smart_context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUnit_Login(t *testing.T) {
	db := test_sqlite.CreateUniqueSqliteInMemoryForUnitTests()

	logger := smart_context.NewSmartContext().WithDB(db)

	hashedPassword, err := helpers.HashPassword("test")
	if err != nil {
		t.Fatalf("Failed to HashPassword: %v", err)
	}

	tests := []struct {
		name       string
		login      string
		password   string
		wantErr    bool
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Login with correct credentials",
			login:      "test",
			password:   "test",
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Create user and login",
			login:      "new_user",
			password:   "wrongpassword",
			wantErr:    true,
			wantStatus: http.StatusOK,
		},
	}

	err = logger.WithTransaction(func(tx *gorm.DB) error {
		// Миграция схемы БД
		if err := tx.AutoMigrate(&model.User{}); err != nil {
			return fmt.Errorf("Failed to migrate database: %v", err)
		}

		err = tx.Create(&model.User{
			ID:           helpers.GenerateUUID(),
			Username:     "test",
			PasswordHash: hashedPassword,
		}).Error
		if err != nil {
			return fmt.Errorf("Failed to create user: %v", err)
		}

		return nil
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, status := login(logger, tt.login, tt.password)

			assert.Equal(t, tt.wantStatus, status)
		})
	}
	assert.NoError(t, err)
}
