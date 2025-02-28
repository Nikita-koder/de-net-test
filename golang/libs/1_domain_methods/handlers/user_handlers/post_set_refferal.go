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

	err := setReferrer(sctx, userID, requestData.ReferrerID)
	if err != nil {
		http.Error(w, `{"errors": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Referral set successfully"}`))
}

func setReferrer(sctx smart_context.ISmartContext, userID, referrerID string) (err error) {
	err = sctx.WithTransaction(func(tx *gorm.DB) error {
		// Проверяем, существует ли реферер
		var referrer model.User
		if err := tx.First(&referrer, "id = ?", referrerID).Error; err != nil {
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
			ReferrerID: referrerID,
			ReferredID: userID,
			CreatedAt:  helpers.GetCurrentTime(),
		}
		if err := tx.Create(&referral).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = awardReferrerBonus(sctx, referrerID, 10)
	if err != nil {
		return err
	}

	return err
}

func awardReferrerBonus(sctx smart_context.ISmartContext, userID string, bonus int32) error {
	return sctx.WithTransaction(func(tx *gorm.DB) error {
		var existingPoints model.Point
		if err := tx.First(&existingPoints, "user_id = ?", userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				existingPoints = model.Point{UserID: &userID, Balance: bonus}
				if err := tx.Create(&existingPoints).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err := tx.Model(&existingPoints).Update("balance", gorm.Expr("balance + ?", bonus)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
