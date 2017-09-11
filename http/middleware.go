package http

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Adapter func(http.Handler) http.Handler

// Adapt h with all specified adapters.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}

	return h
}

func setHeaders() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			h.ServeHTTP(w, r)
		})
	}
}

func checkHeaders(headers map[string]string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for header, value := range headers {
				h := r.Header.Get(header)
				if h == "" || strings.ToLower(h) != strings.ToLower(value) {
					err := fmt.Sprintf("`%s` header must be `%s`.", header, value)
					log.Error(err)

					w.Header().Set("Accept", "application/json")
					w.Header().Set("Accept-Charset", "utf-8")

					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\" }", err)))
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
