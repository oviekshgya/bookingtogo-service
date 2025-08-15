package middleware

import (
	"any/bookingtogo-service/internal/domain"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func Log(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			capture := &responseCaptureWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(capture, r)
			duration := time.Since(start)
			latency := time.Since(start).Milliseconds()

			// Siapkan log data
			reqLog := domain.RequestLog{
				Timestamp:  start,
				IP:         r.RemoteAddr,
				Method:     r.Method,
				Path:       r.URL.Path,
				QueryParam: r.URL.RawQuery,
				UserAgent:  r.UserAgent(),
				Headers:    toJSON(r.Header),
				Body:       domain.NewSqlNullString(string(bodyBytes)),
				Response:   domain.NewSqlNullString(capture.body.String()),
				StatusCode: capture.statusCode,
				CreatedAt:  time.Now(),
				Latency:    latency,
				UpdatedAt:  time.Now(),
			}

			go func(logEntry domain.RequestLog) {
				if db != nil {
					if err := db.Create(&logEntry).Error; err != nil {
						log.Println("Failed to insert request log:", err)
					}
				}
			}(reqLog)

			log.Printf("[%s] %s %s %d %s\n", r.RemoteAddr, r.Method, r.URL.Path, capture.statusCode, duration)
		})
	}

}

func toJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
