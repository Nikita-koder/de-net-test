package user_handlers

import (
	"de-net/libs/1_domain_methods/helpers"
	"de-net/libs/2_generated_models/model"
	"de-net/libs/4_common/smart_context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetReferrerHandler(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var requestData struct {
		ReferrerID string `json:"referrer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"errors": "Неверный запрос"}`, http.StatusBadRequest)
		return
	}

	// Проверка, чтобы пользователь не ввёл свой же ID
	if userID == requestData.ReferrerID {
		http.Error(w, `{"errors": "Нельзя ввести свой же ID"}`, http.StatusBadRequest)
		return
	}

	err := sctx.WithTransaction(func(tx *gorm.DB) error {
		// Проверяем, существует ли реферер
		var referrer model.User
		if err := tx.First(&referrer, "id = ?", requestData.ReferrerID).Error; err != nil {
			return fmt.Errorf("referrer not found")
		}

		// Проверяем, есть ли уже реферальная связь
		var count int64
		tx.Model(&model.Referral{}).Where("referred_id = ?", userID).Count(&count)
		if count > 0 {
			return fmt.Errorf("referral already set")
		}

		// Создаём запись о реферале
		referral := model.Referral{
			ID:         helpers.GenerateUUID(),
			ReferrerID: requestData.ReferrerID,
			ReferredID: userID,
			CreatedAt:  helpers.GetCurrentTime(),
		}
		if err := tx.Create(&referral).Error; err != nil {
			return err
		}

		// Начисляем бонус рефереру (например, 50 поинтов)
		referrerBonus := 50
		if err := tx.Exec("INSERT INTO points (user_id, balance) VALUES (?, ?) ON CONFLICT (user_id) DO UPDATE SET balance = points.balance + ?", requestData.ReferrerID, referrerBonus, referrerBonus).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		http.Error(w, `{"errors": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Referral set successfully"}`))
}
