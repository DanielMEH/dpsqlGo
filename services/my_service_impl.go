package services

import (
	"database/sql"
	"net/http"

	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/ports"
	usecases "github.com/DanielMEH/database/useCases"
)

type MyServiceImpl struct {
	q        *usecases.Query
	queryGet *usecases.QueryGet
	QueryDelete *usecases.QueryDelete
	QueryPut *usecases.QueryPut
	Audits *usecases.QueryAudits

	// Puedes agregar dependencias o campos necesarios
}

func NewMyService() ports.MyService {
	return &MyServiceImpl{
		q:        &usecases.Query{},
		queryGet: &usecases.QueryGet{},
		QueryDelete: &usecases.QueryDelete{},
		QueryPut: &usecases.QueryPut{},
		Audits: &usecases.QueryAudits{},
		// Puedes inicializar dependencias o campos necesarios

	}
}

func (s *MyServiceImpl) Getdb(r *http.Request, dbconn []*sql.DB,log *logger.HeadersRequest) (string, error) {
	return s.queryGet.Getdb(r, dbconn,log)
}

func (s *MyServiceImpl) PostDB(r *http.Request, dbconn []*sql.DB,log *logger.HeadersRequest) (string, error) {
	return s.q.Postdb(r, dbconn,log)
}

func (s *MyServiceImpl) Deletedb(r *http.Request, dbconn []*sql.DB,log *logger.HeadersRequest) (string, error) {
	return s.QueryDelete.Deletedb(r, dbconn,log)
}
func (s *MyServiceImpl) Putdb(r *http.Request, dbconn []*sql.DB,log *logger.HeadersRequest) (string, error) {
	return s.QueryPut.Putdb(r, dbconn,log)
}
