package app

import (
	"github.com/bagasadiii/gofood-clone/handler"
	"github.com/bagasadiii/gofood-clone/middleware"
	"github.com/gorilla/mux"
)

type HandlerDependencies struct {
	UserEndpoint     handler.UserHandlerImpl
	MerchantEndpoint handler.MerchantHandlerImpl
	DriverEndpoint   handler.DriverHandlerImpl
	Middleware       middleware.JWTServiceImpl
}
type Router struct {
	deps HandlerDependencies
}

func NewRouter(deps HandlerDependencies) *Router {
	return &Router{deps: deps}
}

func (ar *Router) Route() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/register", ar.deps.UserEndpoint.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", ar.deps.UserEndpoint.LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/u/{username}", ar.deps.UserEndpoint.GetUserHandler).Methods("GET")

	r.HandleFunc("/api/v1/m/{username}", ar.deps.MerchantEndpoint.GetMerchantHandler).Methods("GET")
	r.HandleFunc("/api/v1/d/{username}", ar.deps.DriverEndpoint.GetDriverHandler).Methods("GET")

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(ar.deps.Middleware.ValidateContext)
	protected.HandleFunc("/m/{username}", ar.deps.MerchantEndpoint.CreateMerchantHandler).Methods("POST")
	protected.HandleFunc("/m/{username}", ar.deps.MerchantEndpoint.UpdateMerchantHandler).Methods("PATCH")

	protected.HandleFunc("/d/{username}", ar.deps.DriverEndpoint.CreateDriverHandler).Methods("POST")
	protected.HandleFunc("/d/{username}", ar.deps.DriverEndpoint.UpdateDriverHandler).Methods("PATCH")
	return r
}
