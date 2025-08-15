package app

import (
	"any/bookingtogo-service/internal/route"
	"any/bookingtogo-service/src/middleware"
	"any/bookingtogo-service/src/pkg"
	"net/http"
	"os"

	muxHandler "github.com/gorilla/handlers"
	muxlogrus "github.com/pytimer/mux-logrus"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Start() {
	if viper.GetString("SERVICE_MODE") == "development" {
		println("Running in development mode")
	}

	r := mux.NewRouter()
	route.AppRoutes(r)

	r.Use(

		muxlogrus.NewLogger(muxlogrus.LogOptions{
			Formatter: &logrus.JSONFormatter{},
		}).Middleware,
	)

	defaultHeader := []middleware.Headers{
		{Key: pkg.H_LANG, Value: "en"},
		{Key: pkg.H_CURRENCY, Value: pkg.CURRENCY_USD},
		{Key: pkg.H_USERID, Value: ""},
		{Key: pkg.H_VISIBILITY, Value: "0"},
		{Key: pkg.H_XTIMEZONE, Value: viper.GetString("TIMEZONE")},
	}

	r.Use(middleware.NewHeaderInformation(defaultHeader...).Middleware)

	err := http.ListenAndServe(":"+viper.GetString("SERVICE_PORT"), muxHandler.CORS(

		muxHandler.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "X-Client", "X-License-ID", "X-License", "X-Device-ID", "X-Session-ID", "token", "X-Api-Key"}),
		muxHandler.AllowedMethods([]string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}),
		muxHandler.AllowedOrigins([]string{"*"}),
		muxHandler.AllowCredentials(),
	)(muxHandler.CompressHandler(muxHandler.LoggingHandler(os.Stdout, r))))
	if err != nil {
		panic(err)
	}
}
