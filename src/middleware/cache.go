package middleware

import (
	"any/bookingtogo-service/src/pkg"
	"any/bookingtogo-service/src/redis"
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type responseCaptureWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseCaptureWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseCaptureWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func generateCacheKey(r *http.Request) string {
	return "request_response:" + r.RemoteAddr + ":" + r.Method + ":" + r.URL.Path
}

func Cache(redisClient *redis.RedisClient) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

			key := generateCacheKey(r)

			var dataCache pkg.RequestLog
			err := redisClient.GetKey(key, &dataCache)
			if err == nil && dataCache.Response != "" {
				log.Println("data from redis")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(dataCache.Response))
				return
			}

			next.ServeHTTP(capture, r)

			log := pkg.RequestLog{
				Timestamp:  time.Now(),
				IP:         r.RemoteAddr,
				Method:     r.Method,
				Path:       r.URL.Path,
				UserAgent:  r.UserAgent(),
				Headers:    r.Header,
				QueryParam: r.URL.RawQuery,
				Body:       string(bodyBytes),
				Response:   capture.body.String(),
				StatusCode: capture.statusCode,
			}

			if err != nil || dataCache.Response == "" {
				_ = redisClient.SetKey(key, log, 5*time.Minute)
			}
		})
	}
}
