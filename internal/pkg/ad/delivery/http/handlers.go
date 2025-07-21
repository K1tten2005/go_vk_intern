package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/ad"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/jwtUtils"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/sendErr"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/validation"
	"github.com/mailru/easyjson"
	"github.com/satori/uuid"
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

	var req models.AdReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while unmarshaling JSON: %w", err), http.StatusBadRequest)
		sendErr.SendError(w, "incorrect request", http.StatusBadRequest)
		return
	}

	adReq := models.Ad{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       int(req.Price * 100.0),
	}
	if err := validation.ValidateAd(adReq); err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
		sendErr.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	adReq.Sanitize()

	userId, ok := jwtUtils.GetIdFromContext(r.Context())
	if !ok {
		logger.LogHandlerError(loggerVar, errors.New("error while getting user id from context"), http.StatusInternalServerError)
		sendErr.SendError(w, "server error", http.StatusInternalServerError)
		return
	}
	login, ok := jwtUtils.GetLoginFromContext(r.Context())
	if !ok {
		logger.LogHandlerError(loggerVar, errors.New("error while getting user login from context"), http.StatusInternalServerError)
		sendErr.SendError(w, "server error", http.StatusInternalServerError)
		return
	}
	adReq.UserId = userId
	adReq.AuthorLogin = login

	advertisement, err := h.uc.CreateAd(r.Context(), adReq)
	if err != nil {
		switch err {
		case ad.ErrCreatingAd:
			logger.LogHandlerError(loggerVar, err, http.StatusInternalServerError)
			sendErr.SendError(w, err.Error(), http.StatusInternalServerError)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusInternalServerError)
			sendErr.SendError(w, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	adResp := models.AdResp{
		Id:          advertisement.Id,
		Title:       advertisement.Title,
		Description: advertisement.Description,
		Price:       float64(advertisement.Price) / 100.0,
		ImageURL:    advertisement.ImageURL,
		CreatedAt:   advertisement.CreatedAt,
		AuthorLogin: advertisement.AuthorLogin,
	}

	data, err := json.Marshal(adResp)
	if err != nil{
		logger.LogHandlerError(loggerVar, fmt.Errorf("error marshaling JSON: %w", err), http.StatusInternalServerError)
		sendErr.SendError(w, "data error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)
}

func (h *AdHandler) GetAds(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var userId uuid.UUID
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userIdStr, ok := jwtUtils.GetIdFromJWT(token, h.secret)
		if ok {
			var err error
			userId, err = uuid.FromString(userIdStr)
			if err != nil {
				userId = uuid.Nil 
			}
		}
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	sortBy := q.Get("sort_by")
	order := q.Get("order")
	priceMin, _ := strconv.Atoi(q.Get("price_min"))
	priceMin *= 100
	priceMax, _ := strconv.Atoi(q.Get("price_max"))
	priceMax *= 100

	if priceMax <= 0 || priceMax > validation.MaxPrice {
		priceMax = validation.MaxPrice
	}
	if priceMin < 0 {
		priceMin = 0
	}
	if priceMin > validation.MaxPrice {
		priceMin = validation.MaxPrice
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if sortBy != "price" && sortBy != "created_at" {
		sortBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	filter := models.Filter{
		Page:     page,
		Limit:    limit,
		SortBy:   sortBy,
		Order:    order,
		PriceMin: priceMin,
		PriceMax: priceMax,
		UserId:   userId,
	}

	ads, err := h.uc.GetAds(r.Context(), filter)
	if err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusInternalServerError)
		sendErr.SendError(w, "failed to load ads", http.StatusInternalServerError)
		return
	}

	resp := make(models.AdRespList, 0, len(ads))
	for _, ad := range ads {
		isOwner := userId != uuid.Nil && ad.UserId == userId
		resp = append(resp, models.AdResp{
			Id:          ad.Id,
			Title:       ad.Title,
			Description: ad.Description,
			Price:       float64(ad.Price) / 100,
			ImageURL:    ad.ImageURL,
			CreatedAt:   ad.CreatedAt,
			AuthorLogin: ad.AuthorLogin,
			IsOwner:     isOwner,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(resp)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("marshal error: %w", err), http.StatusInternalServerError)
		sendErr.SendError(w, "response error", http.StatusInternalServerError)
	}
	w.Write(data)
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)

}
