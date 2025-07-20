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

	authcheck "github.com/K1tten2005/go_vk_intern/internal/pkg/middleware/auth_check"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/middleware/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/router"
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
	logFile, err := os.OpenFile(os.Getenv("MAIN_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("error opening log file: " + err.Error())
		return
	}
	defer logFile.Close()

	loggerVar := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))

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

	r := router.NewRouter()

	r.Use(logMW)

	r.Handle("POST /signin", http.HandlerFunc(authHandler.SignIn))
	r.Handle("POST /signup", http.HandlerFunc(authHandler.SignUp))

	r.Handle("GET /ad", http.HandlerFunc(adHandler.GetAds))
	r.Group(func(r *router.Router) {
		r.Use(authcheck.AuthMiddleware(loggerVar))

		r.Handle("POST /ad", http.HandlerFunc(adHandler.CreateAd))

	})

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
