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

	err := completeTask(sctx, userID, requestData.TaskID)
	if err != nil {
		http.Error(w, `{"errors": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Task completed"}`))
}

func completeTask(sctx smart_context.ISmartContext, userID, taskID string) (err error) {
	err = sctx.WithTransaction(func(tx *gorm.DB) error {
		// Проверяем, существует ли задание
		var task model.Task
		if err := tx.First(&task, "id = ?", taskID).Error; err != nil {
			return errors.New("task not found")
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = awardReferrerBonus(sctx, userID, 20)
	if err != nil {
		return err
	}

	return err
}
