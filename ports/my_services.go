// ports/my_service.go
package ports

import (
	"database/sql"
	"net/http"

	"github.com/DanielMEH/database/logger"
)

type MyService interface {
	Getdb(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error)
	PostDB(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error)
	Deletedb(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error)
	Putdb(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error)
}
