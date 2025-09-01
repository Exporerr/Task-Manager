package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LoggingMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &StatusRecorder{w, http.StatusOK}
		next.ServeHTTP(rec, r)
		duration := time.Since(start)
		ip := r.RemoteAddr
		path := r.URL.Path
		userAgent := r.Header.Get("User-Agent")
		fmt.Printf("[%s] IP:%s|Метод:%s|Путь:%s|Статус:%d|Время:%v|Агент:%s\n", start.Format("2006-01-02 15:04:05"), ip, r.Method, path, rec.status, duration, userAgent)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *StatusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
