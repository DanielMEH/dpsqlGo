package usecases

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/DanielMEH/database/db"
	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/settings"
	"github.com/gorilla/mux"
)

type QueryGet struct{}

func (a *QueryGet) Getdb(httpRequest *http.Request, dbconns []*sql.DB, log *logger.HeadersRequest) (string, error) {

	var (
		query = "SELECT  "
		shema = ""
		table = ""
	)

	respWhere, err := createColumns(httpRequest)
	query += respWhere

	shema = mux.Vars(httpRequest)["schema"]
	table = mux.Vars(httpRequest)["table"]
	query += shema + "." + table + " "

	whereCo, errWhere := CreateWhere(httpRequest)

	op, errOp := Operations(httpRequest)

	query += whereCo
	query += op

	if err != nil {
		return "", err
	}
	if errOp != nil {
		return "", err
	}
	if errWhere != nil {
		return "", err
	}

	resultRows, err := executeGetDb(query, dbconns)
	if err != nil {
		return "", err
	}
	defer resultRows.Close()

	// Procesar los resultados y devolver la información al cliente
	responseData, err := processRows(resultRows)
	if err != nil {
		return "", err
	}

	return responseData, nil

}
func processRows(rows *sql.Rows) (string, error) {
	// Obtener información sobre las columnas
	columns, err := rows.ColumnTypes()
	if err != nil {
		return "", err
	}

	// Preparar un slice para los valores escaneados
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	var responseData []map[string]interface{}

	// Iterar sobre las filas
	for rows.Next() {
		// Escanear los valores en el slice
		if err := rows.Scan(values...); err != nil {
			return "", err
		}

		// Construir un mapa para cada fila
		rowData := make(map[string]interface{})
		for i, col := range columns {
			columnName := col.Name()
			columnValue := *(values[i].(*interface{}))
			rowData[columnName] = columnValue
		}

		// Agregar el mapa al slice de respuesta
		responseData = append(responseData, rowData)
	}

	// Convertir el slice a formato JSON
	jsonResponse, err := json.MarshalIndent(responseData, "", "    ")
	if err != nil {
		return "", err
	}

	// Devolver el JSON como cadena
	return string(jsonResponse), nil
}
func createColumns(http *http.Request) (string, error) {

	var whereConsult = ""
	var dataQuery = http.URL.Query()

	if dataQuery["columns"] == nil {
		logger.LogsWarning("No se envio el parametro columns")
		whereConsult = "*"

	} else {

		for i, columns := range dataQuery["columns"] {
			if i == 0 {
				whereConsult += columns
			} else {
				whereConsult += ", " + columns
			}
		}
	}
	return whereConsult + " FROM ", nil

}

func executeGetDb(query string, dbconns []*sql.DB) (*sql.Rows, error) {
	dbConfig, _ := settings.LoadConfig()
	makeConn := make([]*sql.DB, 0)
	logger.LogsInfo(fmt.Sprintf("query: %v", query))
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
				return nil, errors.New("failed to execute query on all DB connections")
			}

		}

		rows, err := dbconn.Query(query)

		if err != nil {
			// Loguear el error y continuar con el siguiente DB connection
			logger.LogsError(err)
			return nil, err
		} else {
			logger.Logsdebug("Sucessfull query")
			defer dbconn.Close()
			// Loguear el éxito y devolver los resultados
			return rows, nil
		}
	}

	// Si llegamos aquí, significa que todas las conexiones fallaron
	return nil, fmt.Errorf("failed to execute query on all DB connections")
}

func Operations(http *http.Request) (string, error) {

	var dataQuery = http.URL.Query()
	createOperations := ""
	if dataQuery["$limit"] == nil {
		createOperations = ""
	} else {
		for _, value := range dataQuery["$limit"] {
			createOperations = " " + settings.LIMIT + " " + string(value)
		}
	}
	return createOperations, nil

}

func CreateWhere(http *http.Request) (string, error) {
	createWhere := ""
	var dataQuery = http.URL.Query()
	var conditions map[string]interface{}

	if dataQuery["q"] != nil {

		createWhere += " WHERE "
		data := http.URL.Query()["q"][0]
		if err := json.Unmarshal([]byte(data), &conditions); err != nil {
			return "", err
		}
	}

	var conditionsArray []string
	for columns, valor := range conditions {

		valid := fmt.Sprintf("%s = '%v'", columns, valor)
		conditionsArray = append(conditionsArray, valid)
	}

	createWhere += strings.Join(conditionsArray, settings.AND)

	return createWhere, nil

}
