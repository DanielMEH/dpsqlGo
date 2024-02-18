package usecases

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DanielMEH/database/db"
	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/settings"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Importa el driver PostgreSQL
)

type Query struct{}

func (q *Query) Postdb(httpRequest *http.Request, dbconns []*sql.DB, log *logger.HeadersRequest) (string, error) {

	// Luego, lees el cuerpo de la solicitud
	body, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo de la solicitud:", err)
		return "", err
	}

	// Reinicias el puntero del cuerpo para que pueda ser leído nuevamente
	httpRequest.Body = ioutil.NopCloser(strings.NewReader(string(body)))

	var excecutionQuery string = "INSERT INTO "

	schema := mux.Vars(httpRequest)["schema"]
	table := mux.Vars(httpRequest)["table"]

	// im
	excecutionQuery += schema + "." + table + " ("

	var bodyMap map[string]interface{}
	if len(body) == 0 {
		logger.LogsWarning("JSON input is empty")
		return "", errors.New("JSON input is empty")
	}

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		logger.LogsWarning(err)
		return "", err
	}

	// Obtener los campos del body
	var fields []string
	for k := range bodyMap {
		fields = append(fields, k)
	}

	// Agregar los campos al query
	for i, field := range fields {
		excecutionQuery += field
		if i < len(fields)-1 {
			excecutionQuery += ", "
		}
	}

	excecutionQuery += ") VALUES ("

	// Agregar los valores al query
	for i, field := range fields {
		// Verificar si el valor es una cad	 const (

		var value string
		if str, ok := bodyMap[field].(string); ok {
			value = fmt.Sprintf("'%s'", str)
		} else {
			value = fmt.Sprintf("%v", bodyMap[field])
		}

		excecutionQuery += value

		if i < len(fields)-1 {
			excecutionQuery += ", "
		}
	}
	// 3054545

	excecutionQuery += " )"
	config, err := settings.LoadConfig()

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
				return "", errors.New("no se obtuvieron nuevas conexiones")
			}

		}

		err := dbconn.QueryRow(excecutionQuery).Err()
		if err != nil {
			logger.LogsError(fmt.Sprintf("error execute %s", excecutionQuery))
			return "", err
		} else {
			defer dbconn.Close()
			Auditsdb := &QueryAudits{}
			Auditsdb.Auditsdb(httpRequest, dbconns, log)
			logger.LogsInfo(fmt.Sprintf("ok execute  %s", excecutionQuery))

			// cerrar la conexión:
			defer dbconn.Close()

			return "ok execute ", nil
		}

	}

	return "", nil

}
