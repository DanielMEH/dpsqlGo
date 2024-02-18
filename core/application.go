// core/application.go
package core

import (
	"database/sql"
	"net/http"

	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/ports"
	"github.com/DanielMEH/database/services"
)

type Application struct {
	myService ports.MyService
}

func NewApplication() *Application {
	return &Application{
		myService: services.NewMyService(), // Aquí deberías poder acceder a NewMyService
	}
}
 
func (a *Application) Getdb(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error) {
	return a.myService.Getdb(r, dbconn, log)
}
func (a *Application) PostDB(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error) {
	return a.myService.PostDB(r, dbconn,log)
}
func (a *Application) Deletedb(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error) {
	return a.myService.Deletedb(r, dbconn,log)
}
func (a *Application) Putdb(r *http.Request, dbconn []*sql.DB, log *logger.HeadersRequest) (string, error) {
	return a.myService.Putdb(r, dbconn,log)
}
