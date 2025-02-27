package main

import (
	"context"
	"de-net/libs/1_domain_methods/handlers"
	"de-net/libs/1_domain_methods/handlers/user_handlers"
	"de-net/libs/3_infrastructure/db_manager"
	"de-net/libs/4_common/env_vars"
	"de-net/libs/4_common/middleware"
	"de-net/libs/4_common/shutdown"
	"de-net/libs/4_common/smart_context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	env_vars.LoadEnvVars() // load env vars from .env file if ENV_PATH is specified
	BACKEND_PORT, ok := os.LookupEnv("BACKEND_PORT")
	if !ok {
		BACKEND_PORT = "8080"
	}

	sctx := smart_context.NewSmartContext()

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sctx = sctx.WithContext(ctx)
	wg := &sync.WaitGroup{}
	sctx = sctx.WithWaitGroup(wg)

	location, err := time.LoadLocation("Local")
	if err != nil {
		sctx.Fatalf("Error loading local timezone: %v", err)
	}

	sctx.Infof("Local Timezone: %v", location)
	sctx.Infof("Current Time: %v", time.Now().In(location))

	dbm, err := db_manager.NewDbManager(sctx)
	if err != nil {
		sctx.Fatalf("Error connecting to database: %v", err)
	}
	gormDB := dbm.GetGORM()

	// Получаем *sql.DB из gormDB и настраиваем пул соединений
	sqlDB, err := gormDB.DB()
	if err != nil {
		sctx.Fatalf("Error getting sql.DB from GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	// TODO: Подумать и найти место получше для пула соединений

	sctx = sctx.WithDbManager(dbm)
	sctx = sctx.WithDB(gormDB)

	r := chi.NewRouter()

	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "X-Request-Id", "X-Session-Id", "Apikey", "X-Api-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/api/auth", handlers.LoginHandler(sctx))
	// Роуты
	r.Get("/users/{id}/status", middleware.WithRestApiSmartContext(sctx, user_handlers.UserStatusHandler)) // вся доступная информация о пользователе

	r.Get("/users/leaderboard", middleware.WithRestApiSmartContext(sctx, user_handlers.LeaderboardHandler))                                                          // топ пользователей с самым большим балансом
	r.Post("/users/{id}/task/complete", middleware.WithRestApiSmartContext(sctx, func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {})) // выполнение задания
	r.Post("/users/{id}/referrer", middleware.WithRestApiSmartContext(sctx, func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {}))      // ввод реферального кода (может быть id другого пользователя)

	// r.Get("/api/info", middleware.WithRestApiSmartContext(sctx, handlers.InfoHandler))
	// r.Post("/api/sendCoin", middleware.WithRestApiSmartContext(sctx, handlers.SendCoinHandler))
	// r.Get("/api/buy/{item}", middleware.WithRestApiSmartContext(sctx, handlers.BuyItemHandler))

	sctx.Info("Server listening on port " + BACKEND_PORT)
	webServer := &http.Server{
		Addr:         ":" + BACKEND_PORT,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  240 * time.Second,
	}

	go func() {
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sctx.Fatalf("Server error: %v", err)
		}
	}()

	defer func() {
		err := closeFunc(sctx, webServer)
		if err != nil {
			sctx.Errorf("Error closing service: %v", err)
		}
		sctx.Infof("server: Closed")
	}()

	// тут мы зависнем до получения сигнала на завершение
	osSignal := shutdown.WaitForSignalToShutdown()
	// отменяем контекст и ждем завершения всех запросов
	sctx.Infof("Received signal '%s'. Cancelling context", osSignal.String())
	cancel()
	sctx.Infof("Context cancelled - no new requests will be served. Waiting for all prevously started requests to finish")

	wg.Wait() // тут можем ждать долго - до 2х часов - пока запросы все завершатся. новые запросы уже не будут приниматься!
}

func closeFunc(sctx smart_context.ISmartContext, webServer *http.Server) error {
	// Gracefully shut down the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second) // тут уже надо мало ждать
	defer shutdownCancel()
	if err := webServer.Shutdown(shutdownCtx); err != nil {
		sctx.Errorf("Exiting process: Error shutting down server: %v", err)
	} else {
		sctx.Info("Exiting process: Server shut down")
	}

	sctx.Infof("All requests finished. Closing")
	return nil
}
