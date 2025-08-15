package route

import (
	"any/bookingtogo-service/internal/handler"
	"any/bookingtogo-service/src"
	"any/bookingtogo-service/src/pkg"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Api-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CustomerRouter(r *mux.Router, c handler.CustomerHandler) {
	r.HandleFunc("/", c.CreateCustomer).Methods(http.MethodPost)
	r.HandleFunc("/", c.UpdateCustomer).Methods(http.MethodPut)
	r.HandleFunc("/{id}", c.GetCustomerByID).Methods(http.MethodGet)
	r.HandleFunc("/", c.ListCustomersByNationality).Methods(http.MethodGet)
	r.HandleFunc("/{id}", c.DeleteCustomer).Methods(http.MethodDelete)
}

func NasionalityRouter(r *mux.Router, c handler.NasionalityHandler) {
	r.HandleFunc("/{id}", c.GetNasionalityByID).Methods(http.MethodGet)
	r.HandleFunc("/", c.GetAllNasionalities).Methods(http.MethodGet)
}

func AppRoutes(r *mux.Router) {
	r.StrictSlash(true)
	r.Use(corsMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := pkg.PlugResponse(w)
		_ = res.ReplyCustom(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "invalid request",
		})
	}).Methods(http.MethodGet)
	//
	customer, _ := src.InitializeCustomerController()
	CustomerRouter(r.PathPrefix("/customer").Subrouter(), customer)

	nasionality, _ := src.InitializeNasionalityController()
	NasionalityRouter(r.PathPrefix("/nasionality").Subrouter(), nasionality)

	fmt.Println("Server running on :" + viper.GetString("SERVICE_PORT"))
}
