package user_handlers

import (
	"de-net/libs/1_domain_methods/helpers"
	"de-net/libs/2_generated_models/model"
	"de-net/libs/3_infrastructure/test_sqlite"
	"de-net/libs/4_common/smart_context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetReferrer(t *testing.T) {
	db := test_sqlite.CreateUniqueSqliteInMemoryForUnitTests()
	logger := smart_context.NewSmartContext().WithDB(db)

	// Создаём пользователей
	referrerID := helpers.GenerateUUID()
	referredID := helpers.GenerateUUID()

	db.AutoMigrate(&model.User{}, &model.Referral{}, &model.Point{})
	db.Create(&model.User{ID: referrerID, Username: "referrer"})
	db.Create(&model.User{ID: referredID, Username: "referred"})

	tests := []struct {
		name       string
		userID     string
		referrerID string
		wantErr    bool
	}{
		{
			name:       "Successful referral",
			userID:     referredID,
			referrerID: referrerID,
			wantErr:    false,
		},
		{
			name:       "Cannot refer self",
			userID:     referredID,
			referrerID: referredID,
			wantErr:    true,
		},
		{
			name:       "Referrer not found",
			userID:     referredID,
			referrerID: helpers.GenerateUUID(), // Несуществующий ID
			wantErr:    true,
		},
		{
			name:       "Referral already set",
			userID:     referredID,
			referrerID: referrerID,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setReferrer(logger, tt.userID, tt.referrerID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
