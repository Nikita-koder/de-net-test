package handlers

import (
	"de-net/libs/2_generated_models/model"
	"de-net/libs/3_infrastructure/test_sqlite"
	"de-net/libs/4_common/smart_context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Login(t *testing.T) {
	db := test_sqlite.CreateUniqueSqliteInMemoryForUnitTests()

	defer db.Rollback()

	// Миграция схемы БД
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	logger := smart_context.NewSmartContext().WithDB(db)
	password := "securepassword"

	tests := []struct {
		name       string
		login      string
		password   string
		wantErr    bool
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Create user and login",
			login:      "new_user",
			password:   password,
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantResp:   "eyJ",
		},
		{
			name:       "Login with correct password",
			login:      "new_user",
			password:   password,
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantResp:   "eyJ",
		},
		{
			name:       "Login with incorrect password",
			login:      "new_user",
			password:   "wrongpassword",
			wantErr:    true,
			wantStatus: http.StatusUnauthorized,
			wantResp:   `{"errors": "Неверные учетные данные"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, status := login(logger, tt.login, tt.password)

			assert.Equal(t, tt.wantStatus, status)
			assert.Contains(t, resp, tt.wantResp)
		})
	}
}
