package handlers

import (
	"de-net/libs/1_domain_methods/helpers"
	"de-net/libs/2_generated_models/model"
	"de-net/libs/4_common/smart_context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type LoginRequest struct {
	logger   smart_context.ISmartContext
	Login    string
	Password string
}
type LoginQuery struct {
	r            LoginRequest
	responseChan chan LoginResponse
}
type LoginResponse struct {
	token  string
	status int
	err    error
}
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", errors.New("JWT_SECRET is not set")
	}
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
func login(logger smart_context.ISmartContext, login string, password string) (resp string, status int) {
	var user model.User

	err := logger.WithTransaction(func(tx *gorm.DB) error {
		// Ищем пользователя
		if err := tx.First(&user, "username = ?", login).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Создаём нового пользователя
				hashedPassword, err := helpers.HashPassword(password)
				if err != nil {
					return err
				}

				user = model.User{
					ID:           helpers.GenerateUUID(),
					Username:     login,
					PasswordHash: hashedPassword,
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return `{"errors": "Ошибка при создании пользователя"}`, http.StatusInternalServerError
		}
		return `{"errors": "Ошибка базы данных"}`, http.StatusInternalServerError
	}

	// Проверяем пароль
	if !helpers.CheckPassword(user.PasswordHash, password) {
		return `{"errors": "Неверные учетные данные"}`, http.StatusUnauthorized
	}

	// Генерируем JWT
	token, err := GenerateToken(user.ID)
	if err != nil {
		return `{"errors": "Ошибка генерации токена"}`, http.StatusInternalServerError
	}
	return token, http.StatusOK
}

func LoginHandler(sctx smart_context.ISmartContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Login    string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, `{"errors": "Неверный запрос"}`, http.StatusBadRequest)
			return
		}
		response, status := login(sctx, requestData.Login, requestData.Password)
		// Возвращаем HTTP-ответ
		if status != http.StatusOK {
			http.Error(w, response, status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"token": response}); err != nil {
			http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}

	}
}
