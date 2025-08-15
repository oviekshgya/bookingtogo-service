package handler

import (
	"any/bookingtogo-service/internal/domain"
	"any/bookingtogo-service/internal/service"
	"any/bookingtogo-service/src/pkg"
	"any/bookingtogo-service/src/redis"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CustomerHandlerImpl struct {
	service service.CustomerService
	DB      *gorm.DB
	Redis   *redis.RedisClient
}

type CustomerHandler interface {
	CreateCustomer(w http.ResponseWriter, r *http.Request)
	UpdateCustomer(w http.ResponseWriter, r *http.Request)
	DeleteCustomer(w http.ResponseWriter, r *http.Request)
	GetCustomerByID(w http.ResponseWriter, r *http.Request)
	ListCustomersByNationality(w http.ResponseWriter, r *http.Request)
	ListAllCustomers(w http.ResponseWriter, r *http.Request)
}

func NewCustomerHandler(cfg GlobalConfig) CustomerHandler {
	return &CustomerHandlerImpl{service: cfg.ServiceCustomer(), Redis: cfg.GetConnectionRedis(), DB: cfg.GetConnectionDB()}
}

func (c *CustomerHandlerImpl) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)
	req := pkg.PlugRequest(r, w)
	pReq, errParse := pkg.ParseTo[domain.Customer](req)
	if errParse != nil {
		_ = res.ReplyCustom(http.StatusMethodNotAllowed, map[string]interface{}{
			"status":  "error",
			"message": fmt.Sprintf("error: %s your input: %d", errParse.Error(), pReq.NationalityID),
		})
		return
	}

	result, err := c.service.Create(&pReq)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_ = res.Reply(http.StatusCreated, "00", "01", "created", result)
	return
}

func (c *CustomerHandlerImpl) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)
	req := pkg.PlugRequest(r, w)
	pReq, errParse := pkg.ParseTo[domain.Customer](req)
	if errParse != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "invalid request",
		})
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	pReq.ID = id

	result, err := c.service.Update(&pReq)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_ = res.Reply(http.StatusOK, "00", "02", "updated", result)
}

// DELETE
func (c *CustomerHandlerImpl) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "invalid id",
		})
		return
	}

	err = c.service.Delete(uint(id))
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_ = res.Reply(http.StatusOK, "00", "03", "deleted", nil)
}

// GET BY ID
func (c *CustomerHandlerImpl) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "invalid id",
		})
		return
	}

	result, err := c.service.GetById(uint(id))
	if err != nil {
		_ = res.ReplyCustom(http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_ = res.Reply(http.StatusOK, "00", "04", "success", result)
}

// LIST BY NATIONALITY
func (c *CustomerHandlerImpl) ListCustomersByNationality(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)
	vars := mux.Vars(r)
	natIDStr := vars["nationality_id"]
	natID, err := strconv.ParseUint(natIDStr, 10, 64)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "invalid nationality id",
		})
		return
	}

	result, err := c.service.ListByNationalityID(uint(natID))
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_ = res.Reply(http.StatusOK, "00", "05", "success", result)
}

func (h *CustomerHandlerImpl) ListAllCustomers(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)

	// Ambil query param page dan size
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	// Set default jika tidak ada query param
	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	// Panggil service
	data, err := h.service.ListAll(page, pageSize)
	if err != nil {
		_ = res.ReplyCustom(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	_ = res.Reply(http.StatusOK, "00", "001", "ok", data)
}
