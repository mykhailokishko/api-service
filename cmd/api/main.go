package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"card-validator/internal/domain"
	"card-validator/internal/service"
	transporthttp "card-validator/internal/transport/http"
)

func main() {
	// 1. Ініціалізація структурованого логера JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Читання порту із середовища (за замовчуванням 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// 3. Composition Root: Dependency Injection
	clock := domain.NewRealClock()
	validator := service.NewValidator(clock)
	handler := transporthttp.NewCardHandler(validator)

	// 4. Налаштування роутера
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Health check ендпоінт для оркестраторів (Kubernetes / Docker)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// 5. Конфігурація HTTP сервера з таймаутами (захист від Slowloris атак)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 6. Запуск сервера у фоновій горутині
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("Starting HTTP server", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 7. Очікування сигналів ОС для Graceful Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("Server failed to start", slog.String("error", err.Error()))
		os.Exit(1)

	case sig := <-shutdown:
		logger.Info("Shutdown signal received", slog.String("signal", sig.String()))

		// Даємо до 5 секунд на завершення активних HTTP-запитів
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Graceful shutdown failed, forcing close", slog.String("error", err.Error()))
			_ = server.Close()
			os.Exit(1)
		}
		logger.Info("Server stopped gracefully")
	}
}
