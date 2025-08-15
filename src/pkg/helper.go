package pkg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

type HeaderCtx map[string]string

var HeaderCtxKey = struct{ Name string }{Name: "header_ctx_key"}

type DefaultHeader struct {
	UserId          *uint   `json:"user_id"`
	Currency        *string `json:"currency"`
	Language        *string `json:"language"`
	Timezone        *string `json:"timezone"`
	VisibilityLevel *int8   `json:"visibility_level"`
	FeatureLevel    *int    `json:"feature_level"`
}

type RequestLog struct {
	Timestamp  time.Time   `json:"timestamp"`
	IP         string      `json:"ip"`
	Method     string      `json:"method"`
	Path       string      `json:"path"`
	UserAgent  string      `json:"user_agent"`
	Headers    http.Header `json:"headers"`
	QueryParam string      `json:"query_param"`
	Body       string      `json:"body,omitempty"`
	Response   string      `json:"response,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
}

type responseCaptureWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseCaptureWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseCaptureWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func GetClientRequest(r *http.Request) *RequestLog {
	return &RequestLog{
		Timestamp:  time.Now(),
		IP:         r.RemoteAddr,
		Method:     r.Method,
		Path:       r.URL.Path,
		UserAgent:  r.UserAgent(),
		Headers:    r.Header,
		QueryParam: r.URL.RawQuery,
	}
}

func ConvertClientRequestToMd5(log *RequestLog) string {
	data, _ := json.Marshal(log)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
