package user_handlers

import (
	"de-net/libs/4_common/smart_context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserStatusResponse struct {
	UserID  *string `gorm:"user_id"`
	Name    string  `gorm:"username"`
	Balance int     `gorm:"balance"`
}

func UserStatusHandler(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")
	var response UserStatusResponse
	if err := sctx.GetDB().Select("users.id as user_id, username, balance").Joins("JOIN points ON points.user_id = users.id").Where("user_id = ?", userID).Scan(&response).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
}
