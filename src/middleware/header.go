package middleware

import (
	"any/bookingtogo-service/src/pkg"
	"context"
	"fmt"
	"net/http"
)

type Headers struct {
	Key   string
	Value string
}

type HeaderInfoMiddleware struct {
	listHeader []Headers
}

func NewHeaderInformation(headers ...Headers) *HeaderInfoMiddleware {
	if len(headers) == 0 {
		return &HeaderInfoMiddleware{}
	}
	return &HeaderInfoMiddleware{
		listHeader: headers,
	}
}

func (hd *HeaderInfoMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := pkg.TouchRequest(r, w)

		textToLog := []string{}
		hCtx := pkg.HeaderCtx{}
		for _, header := range hd.listHeader {
			if !req.HasHeader(header.Key) {
				r.Header.Set(header.Key, header.Value)
			}
			hCtx[header.Key] = req.Header(header.Key)

			textToLog = append(textToLog, fmt.Sprintf("%s: %s", header.Key, req.Header(header.Key)))
		}

		ctx := context.WithValue(r.Context(), pkg.HeaderCtxKey, hCtx)

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
