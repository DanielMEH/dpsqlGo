package usecases

import (
	"database/sql"
	"errors"
	"net/http"
	"github.com/DanielMEH/database/logger"
	"github.com/gorilla/mux"
)

type QueryDelete struct{}

func (a *QueryDelete) Deletedb(r *http.Request, dbconns []*sql.DB,log *logger.HeadersRequest) (string, error) {

	schema := mux.Vars(r)["schema"]
	table := mux.Vars(r)["table"]
	query := "DELETE FROM "
	query += schema + "." + table

	queryParams := r.URL.Query()

	if len(queryParams) == 0 {
		return "", errors.New("you must send a query")
	}
	if queryParams["q"] == nil {
		return "", errors.New("you must send a query")
	}

	wheraClausura, err := CreateWhere(r)

	if err != nil {
		return "", err
	}

	query += wheraClausura

	// Ejecuta la query

	_, errQuery := ExecuteDb(query, dbconns)

	if errQuery != nil {
		return "", errQuery
	} else {
		return "ok", nil
	}

}
