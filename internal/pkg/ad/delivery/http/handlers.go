package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/ad"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/jwtUtils"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/send_err"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/validation"
	"github.com/mailru/easyjson"
)

type AdHandler struct {
	uc     ad.AdUsecase
	secret string
}

func CreateAdHandler(uc ad.AdUsecase) *AdHandler {
	return &AdHandler{uc: uc, secret: os.Getenv("JWT_SECRET")}
}

func (h *AdHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.Ad
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while unmarshaling JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "incorrect request", http.StatusBadRequest)
		return
	}

	err := validation.ValidateAd(req)
	if err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
		send_err.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Sanitize()

	userId, ok := jwtUtils.GetIdFromContext(r.Context())
	if !ok {
		logger.LogHandlerError(loggerVar, errors.New("error while getting userId from context"), http.StatusInternalServerError)
		send_err.SendError(w, "server error", http.StatusInternalServerError)
		return
	}

	advertisement, err := h.uc.CreateAd(r.Context(), req, userId)
	if err != nil {
		switch err {
		case ad.ErrCreatingAd:
			logger.LogHandlerError(loggerVar, err, http.StatusInternalServerError)
			send_err.SendError(w, err.Error(), http.StatusInternalServerError)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusInternalServerError)
			send_err.SendError(w, "unknown error", http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	adResp := models.AdResp {
		Id: advertisement.Id,
		Title: advertisement.Title,
		Description: advertisement.Description,
		Price: float64(advertisement.Price) / 100.0,
		ImageURL: advertisement.ImageURL,
		CreatedAt: advertisement.CreatedAt,
		IsOwner: advertisement.IsOwner,
	}
	if _, err := easyjson.MarshalToWriter(adResp, w); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error marshaling JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "data error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}
