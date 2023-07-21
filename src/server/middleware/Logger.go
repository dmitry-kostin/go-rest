package middleware

import (
	"fmt"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/urfave/negroni"
	"net/http"
	"time"
)

type Logger struct {
	*pkg.Logger
}

func NewLogger(logger *pkg.Logger) *Logger {
	return &Logger{logger}
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

func (s *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	s.Debug().Msgf("[request] %v", &logEntry{
		Status:   rw.(negroni.ResponseWriter).Status(),
		Duration: time.Since(start),
		Hostname: r.Host,
		Method:   r.Method,
		Path:     r.URL.Path,
	})
}
