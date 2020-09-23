package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const defaultPort = "8081"

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	return &Server{Router: mux.NewRouter()}
}

func (s *Server) Run(port string) error {
	port = configPort(port)

	log.Printf("Listening on port %s", port)

	return http.ListenAndServe(port, s.Router)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (s *Server) Wrap(method string, pattern string, handler HandlerFunc) {
	wrapH := func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		handleError(w, err)
	}

	s.Router.HandleFunc(pattern, wrapH).Methods(method)
}

func configPort(port string) string {
	if port == "" {
		port = defaultPort
	}

	if string(port[0]) != ":" {
		port = fmt.Sprintf(":%s", port)
	}

	return port
}
