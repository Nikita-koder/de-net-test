package user_handlers

import (
	"de-net/libs/4_common/smart_context"
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LeaderboardResponse struct {
	UserID  string `json:"user_id"`
	Name    string `json:"username"`
	Balance int32  `json:"balance"`
}

func LeaderboardHandler(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	// Получаем лимит из параметров запроса (по умолчанию 10)
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var leaderboard []LeaderboardResponse

	err := sctx.WithTransaction(func(tx *gorm.DB) error {
		return tx.Table("users").
			Select("users.id as user_id, users.username, points.balance").
			Joins("JOIN points ON points.user_id = users.id").
			Order("points.balance DESC").
			Limit(limit).
			Find(&leaderboard).Error
	})

	if err != nil {
		http.Error(w, `{"errors": "Ошибка базы данных"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(leaderboard); err != nil {
		http.Error(w, `{"errors": "Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		return
	}
}
