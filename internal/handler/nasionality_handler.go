package handler

import (
	"any/bookingtogo-service/internal/service"
	"any/bookingtogo-service/src/pkg"
	"any/bookingtogo-service/src/redis"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type NasionalityHandlerImpl struct {
	service service.NasionalityService
	Redis   *redis.RedisClient
	DB      *gorm.DB
}

type NasionalityHandler interface {
	GetNasionalityByID(w http.ResponseWriter, r *http.Request)
	GetAllNasionalities(w http.ResponseWriter, r *http.Request)
}

func NewNasionalityHandler(cfg GlobalConfig) NasionalityHandler {
	return &NasionalityHandlerImpl{service: cfg.ServiceNas(), Redis: cfg.GetConnectionRedis(), DB: cfg.GetConnectionDB()}
}

// GET BY ID
func (h *NasionalityHandlerImpl) GetNasionalityByID(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "invalid id",
		})
		return
	}

	result, err := h.service.GetById(id)
	if err != nil {
		_ = res.ReplyCustom(http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_ = res.Reply(http.StatusOK, "00", "04", "success", result)
}

func (h *NasionalityHandlerImpl) GetAllNasionalities(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)

	result, err := h.service.GetAll()
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	_ = res.Reply(http.StatusOK, "00", "05", "success", result)
	return
}
