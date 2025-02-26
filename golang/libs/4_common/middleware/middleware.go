package middleware

import (
	"de-net/libs/4_common/smart_context"
	"errors"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - структура для хранения данных токена
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type SmartHandlerFunc func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request)

func WithRestApiSmartContext(
	externalSctx smart_context.ISmartContext,
	smartHandler SmartHandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WithRecoverer(
			WithWaitGroup(
				WithSmartContext(
					externalSctx,
					smartHandler,
				),
			),
		)(externalSctx, w, r)
	}
}

func WithSmartContext(
	externalSctx smart_context.ISmartContext,
	handler SmartHandlerFunc,
) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecret, ok := os.LookupEnv("JWT_SECRET")
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authSctx := externalSctx.WithField("UserID", token.Claims.(*Claims).UserID)
		handler(authSctx, w, r)
	}
}

func WithWaitGroup(
	handler SmartHandlerFunc,
) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		serverCtx := sctx.GetContext()
		if serverCtx != nil && serverCtx.Err() != nil {
			sctx.Warnf("Server context is closed: %v. Cannot run request", serverCtx.Err())
			http.Error(w, "Server context is closed", http.StatusServiceUnavailable)
			return
		}

		wg := sctx.GetWaitGroup()
		if wg != nil {
			wg.Add(1)
			defer wg.Done()
		}

		handler(sctx, w, r)
	}
}
func WithRecoverer(
	smartHandler SmartHandlerFunc,
) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				stack := debug.Stack()

				sctx.Errorf("RECOVERED FROM UNHANDLED PANIC: %v\nSTACK: %s", panicMessage, stack)
			}
		}()

		smartHandler(sctx, w, r)
	}
}
