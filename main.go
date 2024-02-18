package main

import (
	"net/http"

	"github.com/DanielMEH/database/adapters/httpadapter"
	"github.com/DanielMEH/database/core"
	"github.com/DanielMEH/database/db"
	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/settings"
	"github.com/gorilla/mux"
)

func main() {

	var log = logger.HeadersRequest{}

	config, err := settings.LoadConfig()

	router := mux.NewRouter()

	dbconn, errdb := db.GetDbFromConfig(*config)

	if errdb != nil && dbconn == nil {
		logger.LogsInfo(map[string]interface{}{"error": errdb})
	}
	app := core.NewApplication()

	// Configurar el adaptador HTTP
	httpAdapter := httpadapter.NewHTTPAdapter(app, dbconn, &log)
	// Configurar las rutas
	router.HandleFunc("/{schema:[a-zA-Z_]+}/{table:[a-zA-Z_]+}", httpadapter.MakeCallback(httpAdapter.Getdb)).Methods(http.MethodGet)
	router.HandleFunc("/{schema:[a-zA-Z_]+}/{table:[a-zA-Z_]+}", httpadapter.MakeCallback(httpAdapter.PostDB)).Methods(http.MethodPost)
	router.HandleFunc("/{schema:[a-zA-Z_]+}/{table:[a-zA-Z_]+}", httpadapter.MakeCallback(httpAdapter.PutDB)).Methods(http.MethodPut)
	router.HandleFunc("/{schema:[a-zA-Z_]+}/{table:[a-zA-Z_]+}", httpadapter.MakeCallback(httpAdapter.DeleteDB)).Methods(http.MethodDelete)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("404"))
	})

	if err != nil {
		panic(err)
	}
	// listening on server
	logger.LogsInfo("listening on server http://" + config.ConfigDb.Host + "" + config.ConfigDb.Port + " ")

	http.Handle("/", router)
	http.ListenAndServe(config.ConfigDb.Port, nil)

}
