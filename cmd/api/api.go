package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AdvenAdam/ecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api").Subrouter()
	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subRouter)

	log.Println("Listening on port", s.addr)

	return http.ListenAndServe(s.addr, router)
}
