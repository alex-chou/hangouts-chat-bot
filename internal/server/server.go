package server

import (
	"net/http"
	"time"

	"github.com/alex-chou/chat-bot/internal/backend"
)

const (
	// timeoutDuration is the default timeout duration for every handler.
	timeoutDuration = 5 * time.Second
	// timeoutMessage is the default timeout message for every handler.
	timeoutMessage = "{\"error\": \"endpoint timeout\"}"
)

// Server is a chatbot server.
type Server struct {
	*http.ServeMux
	backend backend.Backend
	Token   string
}

// NewServer creates a new server.
func NewServer(backend backend.Backend) *Server {
	server := &Server{
		ServeMux: http.NewServeMux(),
		backend:  backend,
	}

	server.Handle("/health", NewWrappedHandler(http.MethodGet, server.Health))
	server.Handle("/supamonkey_slack", NewWrappedHandler(http.MethodPost, server.Slack))

	return server
}

// NewWrappedHandler wraps the handler and passes it to the embedded ServeMux.
func NewWrappedHandler(method string, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	wrappedHandler := http.Handler(http.HandlerFunc(handler))
	wrappedHandler = NewMethodHandler(method, wrappedHandler)
	wrappedHandler = http.TimeoutHandler(wrappedHandler, timeoutDuration, timeoutMessage)
	return wrappedHandler
}

// NewMethodHandler invokes the input handler iff the request method matches the
// input method.
func NewMethodHandler(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handler.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})

}
