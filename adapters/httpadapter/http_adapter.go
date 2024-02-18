// adapters/httpadapter/http_adapter.go
package httpadapter

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/DanielMEH/database/core"
	"github.com/DanielMEH/database/logger"
)

type HTTPAdapter struct {
	app     *core.Application
	dbconns []*sql.DB
	log     *logger.HeadersRequest
}

// definir type como funcion para pasr los los por parametros logs logger.Logger

func NewHTTPAdapter(app *core.Application, dbconns []*sql.DB, log *logger.HeadersRequest) *HTTPAdapter {

	return &HTTPAdapter{app: app, dbconns: dbconns, log: log}
}

func (a *HTTPAdapter) Getdb(w http.ResponseWriter, r *http.Request) {

	log := a.MiddlewareHttp(r)

	data, err := a.app.Getdb(r, a.dbconns, log)
	if err != nil {
		errorMessage := fmt.Sprintf("::%v", err)
		logger.LogsError(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)

		return
	}

	// Enviar respuesta al cliente
	w.Write([]byte(data))
}

func (a *HTTPAdapter) PostDB(w http.ResponseWriter, r *http.Request) {

	log := a.MiddlewareHttp(r)

	data, err := a.app.PostDB(r, a.dbconns, log)
	// intanceLog(w, r)

	if err != nil {
		errorMessage := fmt.Sprintf("::%v", err)
		logger.LogsError(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)

		return
	}

	logger.Logsdebug(data)
	// Enviar respuesta al cliente
	w.Write([]byte(data))
}

func (a *HTTPAdapter) PutDB(w http.ResponseWriter, r *http.Request) {

	log := a.MiddlewareHttp(r)
	data, err := a.app.Putdb(r, a.dbconns, log)

	if err != nil {
		errorMessage := fmt.Sprintf("::%v", err)
		logger.LogsError(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)

		return
	}

	// Enviar respuesta al cliente
	w.Write([]byte(data))

}
func (a *HTTPAdapter) DeleteDB(w http.ResponseWriter, r *http.Request) {

	log := a.MiddlewareHttp(r)
	data, err := a.app.Deletedb(r, a.dbconns, log)
	if err != nil {
		errorMessage := fmt.Sprintf("::%v", err)
		logger.LogsError(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)

		return
	}

	// Enviar respuesta al cliente
	w.Write([]byte(data))

}

func (a *HTTPAdapter) MiddlewareHttp(r *http.Request) *logger.HeadersRequest {

	logger.LogsInfo(r.Header)

	headersRequest := &logger.HeadersRequest{
		Authorization: r.Header.Get("Authorization"),
		Actions:       r.Header.Get("Action"),
		LogSecId:      r.Header.Get("logId"),
		Ip:            r.Header.Get("Ip"),
		Username:      r.Header.Get("Username"),
	}
	logger.BuildTransformData(headersRequest)
	return headersRequest
}

func MakeCallback(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
