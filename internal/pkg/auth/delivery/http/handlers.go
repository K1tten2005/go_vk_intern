package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/sendErr"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/validation"
	"github.com/mailru/easyjson"
)

type AuthHandler struct {
	uc     auth.AuthUsecase
	secret string
}

func CreateAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc, secret: os.Getenv("JWT_SECRET")}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.UserReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while unmarshaling JSON: %w", err), http.StatusBadRequest)
		sendErr.SendError(w, "incorrect request", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	if !validation.ValidPassword(req.Password) {
		logger.LogHandlerError(loggerVar, auth.ErrInvalidPassword, http.StatusBadRequest)
		sendErr.SendError(w, auth.ErrInvalidPassword.Error(), http.StatusBadRequest)
		return
	}

	if !validation.ValidLogin(req.Login) {
		logger.LogHandlerError(loggerVar, auth.ErrInvalidLogin, http.StatusBadRequest)
		sendErr.SendError(w, auth.ErrInvalidLogin.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.uc.SignIn(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrUserNotFound, auth.ErrInvalidCredentials:
			logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
			sendErr.SendError(w, err.Error(), http.StatusBadRequest)
		case auth.ErrGeneratingToken:
			logger.LogHandlerError(loggerVar, err, http.StatusInternalServerError)
			sendErr.SendError(w, err.Error(), http.StatusInternalServerError)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusInternalServerError)
			sendErr.SendError(w, "unknown error", http.StatusInternalServerError)
		}
		return
	}
	user.Token = token

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(user)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error marshaling JSON: %w", err), http.StatusInternalServerError)
		sendErr.SendError(w, "data error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.UserReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while unmarshalling JSON: %w", err), http.StatusBadRequest)
		sendErr.SendError(w, "incorrect request", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	if !validation.ValidPassword(req.Password) {
		logger.LogHandlerError(loggerVar, auth.ErrInvalidPassword, http.StatusBadRequest)
		sendErr.SendError(w, auth.ErrInvalidPassword.Error(), http.StatusBadRequest)
		return
	}

	if !validation.ValidLogin(req.Login) {
		logger.LogHandlerError(loggerVar, auth.ErrInvalidLogin, http.StatusBadRequest)
		sendErr.SendError(w, auth.ErrInvalidLogin.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.uc.SignUp(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrInvalidPassword:
			logger.LogHandlerError(loggerVar, fmt.Errorf("invalid password: %w", err), http.StatusBadRequest)
			sendErr.SendError(w, err.Error(), http.StatusBadRequest)
		case auth.ErrInvalidLogin:
			logger.LogHandlerError(loggerVar, fmt.Errorf("invalid login: %w", err), http.StatusBadRequest)
			sendErr.SendError(w, err.Error(), http.StatusBadRequest)
		case auth.ErrUserAlreadyExists:
			logger.LogHandlerError(loggerVar, err, http.StatusConflict)
			sendErr.SendError(w, err.Error(), http.StatusConflict)
		case auth.ErrCreatingUser:
			logger.LogHandlerError(loggerVar, err, http.StatusInternalServerError)
			sendErr.SendError(w, err.Error(), http.StatusInternalServerError)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusInternalServerError)
			sendErr.SendError(w, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	user.Token = token

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(user)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error marshaling JSON: %w", err), http.StatusInternalServerError)
		sendErr.SendError(w, "data error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)
}
