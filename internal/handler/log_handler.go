package handler

import (
	"any/bookingtogo-service/internal/service"
	"any/bookingtogo-service/src/pkg"
	"net/http"
	"strconv"
)

type RequestLogHandler interface {
	ListAllLogs(w http.ResponseWriter, r *http.Request)
}

type RequestLogHandlerImpl struct {
	service service.RequestLogService
}

func NewRequestLogHandler(cfg GlobalConfig) RequestLogHandler {
	return &RequestLogHandlerImpl{
		service: cfg.LogService(),
	}
}

func (h *RequestLogHandlerImpl) ListAllLogs(w http.ResponseWriter, r *http.Request) {
	res := pkg.PlugResponse(w)

	// Ambil query param page dan size
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

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

	data, err := h.service.ListAll(page, pageSize)
	if err != nil {
		_ = res.ReplyCustom(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Custom response agar mudah ditangkap Laravel
	response := map[string]interface{}{
		"status":     "success",
		"page":       data.Page,
		"page_size":  data.PageSize,
		"total_page": data.TotalPage,
		"total_data": data.TotalData,
		"logs": func() []map[string]interface{} {
			list := []map[string]interface{}{}
			for _, log := range data.Logs {
				list = append(list, map[string]interface{}{
					"id":         log.ID,
					"timestamp":  log.Timestamp,
					"ip":         log.IP,
					"method":     log.Method,
					"path":       log.Path,
					"status":     log.StatusCode,
					"latency_ms": log.Latency,
				})
			}
			return list
		}(),
	}

	_ = res.Reply(http.StatusOK, "00", "001", "ok", response)
}
