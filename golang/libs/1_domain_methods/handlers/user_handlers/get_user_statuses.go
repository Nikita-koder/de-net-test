package user_handlers

import (
	"de-net/libs/4_common/smart_context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type UserStatusResponse struct {
	UserID  *string `gorm:"user_id"`
	Name    string  `gorm:"username"`
	Balance int     `gorm:"balance"`
}

func UserStatusHandler(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")
	var response UserStatusResponse

	err := sctx.WithTransaction(func(tx *gorm.DB) error {
		return tx.Table("users").
			Select("users.id as user_id, users.username, points.balance").
			Joins("JOIN points ON points.user_id = users.id").
			Where("users.id = ?", userID).
			Find(&response).Error
	})

	if err != nil {
		http.Error(w, `{"errors": "User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"errors": "Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		return
	}
}
