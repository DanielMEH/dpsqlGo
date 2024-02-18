package usecases

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DanielMEH/database/db"
	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/settings"
	"github.com/gorilla/mux"
)

type QueryPut struct{}

func (a *QueryPut) Putdb(httpRequest *http.Request, dbconns []*sql.DB, log *logger.HeadersRequest) (string, error) {

	var query = "UPDATE "

	schema := mux.Vars(httpRequest)["schema"]
	table := mux.Vars(httpRequest)["table"]

	query += schema + "." + table

	strResponse, err := createSet(httpRequest)

	if err != nil {
		return "", err
	}
	strWhere, errw := CreateWhere(httpRequest)

	if errw != nil {
		return "", errw
	}

	query += strResponse
	query += strWhere

	_, err = ExecuteDb(query, dbconns)

	if err != nil {
		return "", err
	}
	return "ok execute", nil
}

func createSet(http *http.Request) (string, error) {

	var setClousure = " "
	body, errs := ioutil.ReadAll(http.Body)

	if errs != nil {
		return "", errs
	}
	http.Body = ioutil.NopCloser(http.Body)

	if len(body) == 0 {
		return "", errors.New("JSON input is empty")
	}
	// var bodyMap = map[name:duvan passkey:56565656 surname:garzon]
	var bodyMap map[string]interface{}

	err := json.Unmarshal(body, &bodyMap)
	if len(body) == 0 {
		logger.LogsError(err)
		return "", err
	}

	var fields []string

	for k := range bodyMap {
		fields = append(fields, k)

	}

	for i, field := range fields {
		setClousure += field + " = " + "'" + bodyMap[field].(string) + "'"
		if i < len(fields)-1 {
			setClousure += ", "
		}
	}

	if err != nil {
		return "", err

	}

	return " SET " + setClousure, nil

}

func ExecuteDb(query string, dbconns []*sql.DB) (string, error) {
	dbConfig, _ := settings.LoadConfig()
	makeConn := make([]*sql.DB, 0)
	for _, dbconn := range dbconns {

		// Verificar si la conexión está cerrada
		if err := dbconn.Ping(); err != nil {
			// Si está cerrada, intenta abrir nuevamente la conexión
			newDB, _ := db.GetDbFromConfig(*dbConfig)

			if len(newDB) > 0 {
				makeConn = append(makeConn, newDB...)
				dbconn = makeConn[0]
			} else {
				// Si no se obtuvieron nuevas conexiones, devolver un error
				return "", errors.New("failed to execute query on all DB connections")
			}

		}

		_, err := dbconn.Exec(query)
		if err != nil {
			logger.LogsError(err)
			return "", err
		} else {
			defer dbconn.Close()
			logger.LogsInfo(fmt.Sprintf("ok execute query: %v", query))
			return "ok execute", nil
		}
	}

	return "ok execute", nil
}
