package authCheck

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/jwtUtils"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/sendErr"
	"github.com/gorilla/mux"
)

func AuthMiddleware(loggerVar *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.LogHandlerError(loggerVar, fmt.Errorf("missing Authorization header"), http.StatusUnauthorized)
				sendErr.SendError(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				logger.LogHandlerError(loggerVar, fmt.Errorf("invalid Authorization header format"), http.StatusUnauthorized)
				sendErr.SendError(w, "invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			secret := os.Getenv("JWT_SECRET")
			id, ok := jwtUtils.GetIdFromJWT(parts[1], secret)
			if !ok {
				logger.LogHandlerError(loggerVar, fmt.Errorf("invalid token"), http.StatusUnauthorized)
				sendErr.SendError(w, "invalid token", http.StatusUnauthorized)
				return
			}

			login, ok := jwtUtils.GetLoginFromJWT(parts[1], secret)
			if !ok {
				logger.LogHandlerError(loggerVar, fmt.Errorf("invalid token"), http.StatusUnauthorized)
				sendErr.SendError(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), jwtUtils.UserIdKey, id)
			ctx = context.WithValue(ctx, jwtUtils.UserLoginKey, login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
