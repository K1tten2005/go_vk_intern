package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/K1tten2005/go_vk_intern/internal/pkg/middleware/authCheck"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/middleware/logger"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	authHandler "github.com/K1tten2005/go_vk_intern/internal/pkg/auth/delivery/http"
	authRepo "github.com/K1tten2005/go_vk_intern/internal/pkg/auth/repo"
	authUsecase "github.com/K1tten2005/go_vk_intern/internal/pkg/auth/usecase"

	adHandler "github.com/K1tten2005/go_vk_intern/internal/pkg/ad/delivery/http"
	adRepo "github.com/K1tten2005/go_vk_intern/internal/pkg/ad/repo"
	adUsecase "github.com/K1tten2005/go_vk_intern/internal/pkg/ad/usecase"
)

func initDB(logger *slog.Logger) (*pgxpool.Pool, error) {
	connStr := os.Getenv("POSTGRES_CONNECTION")

	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Successful connection to PostgreSQL")
	return pool, nil
}

func main() {
	logPath := os.Getenv("MAIN_LOG_FILE")
	logDir := "./logs"

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("failed to create logs directory: " + err.Error())
		return
	}

	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("error opening log file: " + err.Error())
		return
	}
	defer logFile.Close()

	loggerVar := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	pool, err := initDB(loggerVar)
	if err != nil {
		loggerVar.Error("Error while connecting to PostgreSQL: " + err.Error())
		return
	}
	defer pool.Close()

	logMW := logger.CreateLoggerMiddleware(loggerVar)

	authRepo := authRepo.CreateAuthRepo(pool)
	authUsecase := authUsecase.CreateAuthUsecase(authRepo)
	authHandler := authHandler.CreateAuthHandler(authUsecase)

	adRepo := adRepo.CreateAdRepo(pool)
	adUsecase := adUsecase.CreateAdUsecase(adRepo)
	adHandler := adHandler.CreateAdHandler(adUsecase)

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Не найдено", http.StatusNotFound)
	})
	r.Use(logMW)

	publicRoutes := r.PathPrefix("/api").Subrouter()
	publicRoutes.HandleFunc("/signin", authHandler.SignIn).Methods(http.MethodPost)
	publicRoutes.HandleFunc("/signup", authHandler.SignUp).Methods(http.MethodPost)
	publicRoutes.HandleFunc("/ad", adHandler.GetAds).Methods(http.MethodGet)

	protectedRoutes := r.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(authCheck.AuthMiddleware(loggerVar))
	protectedRoutes.HandleFunc("/ad", adHandler.CreateAd).Methods(http.MethodPost)

	srv := http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			loggerVar.Error("Error while starting server: " + err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	loggerVar.Info("Got stop signal")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		loggerVar.Error("Error while stopping server: " + err.Error())
	} else {
		loggerVar.Info("Server successfully stopped")
	}
}
