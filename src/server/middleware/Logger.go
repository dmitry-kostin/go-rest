package middleware

import (
	"fmt"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"net/http"
	"time"
)

type StatusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *StatusCodeResponseWriter) Status() int {
	return s.statusCode
}

func (s *StatusCodeResponseWriter) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

type logEntry struct {
	Status   int           `json:"status"`
	Duration time.Duration `json:"duration"`
	Hostname string        `json:"hostname"`
	Method   string        `json:"method"`
	Path     string        `json:"path"`
}

func (s logEntry) String() string {
	return fmt.Sprintf("%d | \t %s | %s | %s %s", s.Status, s.Duration, s.Hostname, s.Method, s.Path)
}

type Logger struct {
	*pkg.Logger
}

func NewLogger(logger *pkg.Logger) *Logger {
	return &Logger{logger}
}

func (s *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusCodeRW := &StatusCodeResponseWriter{rw, http.StatusOK}
		next.ServeHTTP(statusCodeRW, r)
		s.Debug().Msgf("[request] %v", &logEntry{
			Status:   statusCodeRW.statusCode,
			Duration: time.Since(start),
			Hostname: r.Host,
			Method:   r.Method,
			Path:     r.URL.Path,
		})
	})
}
