package helpers

import (
	"crypto/rand"
	"de-net/libs/4_common/smart_context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateUUID() *string {
	uuid := uuid.NewString()
	return &uuid
}

func GetCurrentTime() *time.Time {
	now := time.Now()
	return &now
}

// Хеширование пароля
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Проверка пароля
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
func GenerateJWT(sctx smart_context.ISmartContext, userId, ip string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userId
	claims["ip"] = ip
	claims["exp"] = time.Now().Add(15 * time.Minute).Unix() // 15 минут

	return token.SignedString([]byte(sctx.GetDbManager().GetJwtSecret()))
}

func GenerateRandomBase64(n int) (string, error) {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func ParseJWT(sctx smart_context.ISmartContext, tokenString string) (string, error) {
	secret := sctx.GetDbManager().GetJwtSecret()
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		if t.Method != jwt.SigningMethodHS512 {
			return nil, errors.New("unexpected signing method (want HS512)")
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userId, _ := claims["user_id"].(string)
	return userId, nil
}
