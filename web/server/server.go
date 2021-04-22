package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const defaultPort = ":8081"

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	return &Server{
		Router: mux.NewRouter(),
	}
}

func (s *Server) Run(port string) error {
	port = setupPort(port)

	s.Wrap(http.MethodGet, "/ping", func(w http.ResponseWriter, r *http.Request) error {
		return RespondJSON(w, "pong", http.StatusOK)
	})

	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(port, s.Router)
}

func setupPort(port string) string {
	if port == "" {
		port = defaultPort
		log.Printf("Defaulting to port %s", port)
	}

	return port
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type Middleware func(h http.HandlerFunc) http.HandlerFunc

func (s *Server) Wrap(method string, pattern string, handler HandlerFunc, mws ...Middleware) {
	s.Router.HandleFunc(pattern, wrapHandler(handlerAdapter(handler), mws...)).Methods(method)
}

func wrapHandler(handler http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
	length := len(mws) - 1
	for mw := length; mw >= 0; mw-- {
		h := mws[mw]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

func handlerAdapter(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		hErr := handleError(w, err)
		_ = RespondJSON(w, hErr, hErr.StatusCode)
	}
}