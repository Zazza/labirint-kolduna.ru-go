package middlewares

import (
	"log"
	"net/http"

	apperrors "gamebook-backend/pkg/errors"
)

type HTTPHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type errorHandler struct {
	next http.Handler
}

func (e *errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v", rec)
			apperrors.WriteHTTPError(w, *apperrors.ErrInternal)
		}
	}()

	e.next.ServeHTTP(w, r)
}

func ErrorHandler(next http.Handler) http.Handler {
	return &errorHandler{next: next}
}

type recoveryHandler struct {
	next http.Handler
}

func (r *recoveryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Recovered-Panic") != "" {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Recovered from panic via recovery middleware: %v", rec)
				apperrors.WriteHTTPError(w, *apperrors.ErrInternal)
			}
		}()
	}
	r.next.ServeHTTP(w, req)
}

func RecoveryHandler(next http.Handler) http.Handler {
	return &recoveryHandler{next: next}
}

type logRequestHandler struct {
	next http.Handler
}

func (l *logRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)
	l.next.ServeHTTP(w, r)
}

func LogRequest(next http.Handler) http.Handler {
	return &logRequestHandler{next: next}
}
