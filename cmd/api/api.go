package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AdvenAdam/go-ecom/service/cart"
	"github.com/AdvenAdam/go-ecom/service/order"
	"github.com/AdvenAdam/go-ecom/service/product"
	"github.com/AdvenAdam/go-ecom/service/user"
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
	// We define a subrouter for all routes that start with "/api/v1". This
	// allows us to have a clear separation of concerns between the API and
	// any other routes that may be defined in the future.
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// We create a new user store using the database connection provided by
	// the user. This allows us to decouple the user store from the database
	// connection and make it easier to test.
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subRouter)

	log.Println("Listening on port", s.addr)

	return http.ListenAndServe(s.addr, router)
}
