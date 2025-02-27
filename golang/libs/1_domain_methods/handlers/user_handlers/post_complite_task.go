package user_handlers

import (
	"de-net/libs/2_generated_models/model"
	"de-net/libs/4_common/smart_context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func CompleteTaskHandler(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var requestData struct {
		TaskID string `json:"task_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"errors": "Неверный запрос"}`, http.StatusBadRequest)
		return
	}

	err := sctx.WithTransaction(func(tx *gorm.DB) error {
		// Проверяем, существует ли задание
		var task model.Task
		if err := tx.First(&task, "id = ?", requestData.TaskID).Error; err != nil {
			return errors.New("task not found")
		}

		// Начисляем поинты пользователю
		if err := tx.Exec("INSERT INTO points (user_id, balance) VALUES (?, ?) ON CONFLICT (user_id) DO UPDATE SET balance = points.balance + ?", userID, task.PointsReward, task.PointsReward).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		http.Error(w, `{"errors": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Task completed"}`))
}
