package route

import (
	"any/bookingtogo-service/internal/handler"
	"any/bookingtogo-service/src"
	"any/bookingtogo-service/src/middleware"
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
	r.HandleFunc("", c.CreateCustomer).Methods(http.MethodPost)
	r.HandleFunc("", c.UpdateCustomer).Methods(http.MethodPut)
	r.HandleFunc("/{id}", c.GetCustomerByID).Methods(http.MethodGet)
	r.HandleFunc("/nasionality", c.ListCustomersByNationality).Methods(http.MethodGet)
	r.HandleFunc("", c.ListAllCustomers).Methods(http.MethodGet)
	r.HandleFunc("/{id}", c.DeleteCustomer).Methods(http.MethodDelete)
}

func LogRouter(r *mux.Router, c handler.RequestLogHandler) {
	r.HandleFunc("", c.ListAllLogs).Methods(http.MethodGet)
}

func NasionalityRouter(r *mux.Router, c handler.NasionalityHandler) {

	log := middleware.Log(GlobalCfg.GetConnectionDB())
	r.Use(log)
	cache := middleware.Cache(GlobalCfg.GetConnectionRedis())
	r.Handle("/{id}", cache(http.HandlerFunc(c.GetNasionalityByID))).Methods(http.MethodGet)
	r.Handle("/", cache(http.HandlerFunc(c.GetAllNasionalities))).Methods(http.MethodGet)
}

var GlobalCfg handler.GlobalConfig

func init() {
	GlobalCfg, _ = src.InitializeGlobalConfigs()
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

	CustomerRouter(r.PathPrefix("/customer").Subrouter(), handler.NewCustomerHandler(GlobalCfg))

	NasionalityRouter(r.PathPrefix("/nasionality").Subrouter(), handler.NewNasionalityHandler(GlobalCfg))

	LogRouter(r.PathPrefix("/log").Subrouter(), handler.NewRequestLogHandler(GlobalCfg))

	fmt.Println("Server running on :" + viper.GetString("SERVICE_PORT"))
}
