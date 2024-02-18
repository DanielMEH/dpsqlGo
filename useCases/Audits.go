package usecases

import (
	"database/sql"
	"encoding/json"
	"errors"

	//"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DanielMEH/database/db"
	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/settings"
	"github.com/gorilla/mux"
)

type QueryAudits struct{}

func (q *QueryAudits) Auditsdb(httpRequest *http.Request, dbconns []*sql.DB, log *logger.HeadersRequest) (string, error) {
	// print httpRequest.Method
	config, _ := settings.LoadConfig()
	makeConn := make([]*sql.DB, 0)

	headersJSON, err := json.Marshal(httpRequest.Header)
	if err != nil {
		// Manejar el error
		fmt.Println("Error al convertir las cabeceras a JSON:", err)
		return "", err
	}
	// Suponiendo que 'body' es un map o una estructura

	body, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo de la solicitud:", err)
		return "", err
	}

	// valid rows

	// 	// Reinicias el puntero del cuerpo para que pueda ser leído nuevamente
	httpRequest.Body = ioutil.NopCloser(strings.NewReader(string(body)))

	var excecutionQuery string = "INSERT INTO "

	schema := mux.Vars(httpRequest)["schema"]
	table := mux.Vars(httpRequest)["table"]

	queryAuditsAll := "SELECT * FROM audits.audittable WHERE SCHEMANAME = '" + strings.ToUpper(schema) + "' AND TABLENAME = '" +
		strings.ToUpper(table) + "'" + " AND STATE = 1"

	rows, errA := ExecuteDbPg(queryAuditsAll, dbconns)
	if errA != nil {
		return "", errA
	}
	defer rows.Close()
	if !rows.Next() {
		logger.LogsError(fmt.Sprintf("Error Audits %v", "No rows found"))
		return "", nil
	}
	// si rows no hay nada

	// 	// im
	var ACTIONQUERY string
	switch httpRequest.Method {
	case "POST":
		ACTIONQUERY = "INSERT"
	case "GET":
		ACTIONQUERY = "SELECT"
	case "PUT":
		ACTIONQUERY = "UPDATE"
	case "DELETE":
		ACTIONQUERY = "DELETE"
	default:
		ACTIONQUERY = "OTHER"
	}

	excecutionQuery += "AUDITS" + "." + "AUDITLOG" + ` (method,ACTIONQUERY, content, username, ip, request, schemaname, tablename)`
	excecutionQuery += " VALUES ('" + httpRequest.Method + "','" + ACTIONQUERY + "','" + string(body) + "','" + log.Username + "','" + log.Ip + "','" + string(headersJSON) + "','" + schema + "','" + table + "')"

	for _, dbconn := range dbconns {
		// Verificar si la conexión está cerrada
		if err := dbconn.Ping(); err != nil {
			// Si está cerrada, intenta abrir nuevamente la conexión
			newDB, _ := db.GetDbFromConfig(*config)

			if len(newDB) > 0 {
				makeConn = append(makeConn, newDB...)
				dbconn = makeConn[0]
			} else {
				// Si no se obtuvieron nuevas conexiones, devolver un error
				return "", errors.New("no se obtuvieron nuevas conexiones")
			}

		}

		err := dbconn.QueryRow(excecutionQuery).Err()
		if err != nil {
			logger.LogsError(fmt.Sprintf("error execute %s ERRORDB: %v", excecutionQuery, err))
			return "", err
		} else {
			defer dbconn.Close()
			logger.Logsdebug(fmt.Sprintf("ok execute  %s", excecutionQuery))
			return "ok execute ", nil
		}

	}

	return "ok", nil

}
func ExecuteDbPg(query string, dbconns []*sql.DB) (*sql.Rows, error) {
	logger.LogsInfo(fmt.Sprintf("query: %v", query))

	config, _ := settings.LoadConfig()
	makeConn := make([]*sql.DB, 0)

	for _, dbconn := range dbconns {

		// Verificar si la conexión está cerrada
		if err := dbconn.Ping(); err != nil {
			// Si está cerrada, intenta abrir nuevamente la conexión
			newDB, _ := db.GetDbFromConfig(*config)

			if len(newDB) > 0 {
				makeConn = append(makeConn, newDB...)
				dbconn = makeConn[0]
			} else {
				// Si no se obtuvieron nuevas conexiones, devolver un error
				return nil, errors.New("no se obtuvieron nuevas conexiones")
			}

		}

		rows, err := dbconn.Query(query)

		if err != nil {
			// Loguear el error y continuar con el siguiente DB connection
			logger.LogsError(err)
			return nil, err
		} else {
			defer dbconn.Close()
			logger.Logsdebug("sucessful Audits")
			// Loguear el éxito y devolver los resultados
			return rows, nil
		}
	}

	// Si llegamos aquí, significa que todas las conexiones fallaron
	return nil, fmt.Errorf("failed to execute query on all DB connections")
}
