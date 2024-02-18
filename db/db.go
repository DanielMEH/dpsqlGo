package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DanielMEH/database/logger"
	"github.com/DanielMEH/database/settings"
	_ "github.com/lib/pq"
)

type DB struct{}

func GetDbFromConfig(config settings.Config) ([]*sql.DB, error) {
	log.Println("Conexiones a iniciar: ", len(config.DbListConn))

	// Crear un slice para almacenar las conexiones abiertas
	var dbConnections []*sql.DB

	// Imprimir settings_Config
	for _, v := range config.DbListConn {
		// Open up our database connections.
		dbconn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", v.Host, v.Port, v.User, v.Password, v.Database))
		if err != nil {
			logger.LogsError("Error al establecer la conexion: " + v.Database + "")
			continue // Saltar a la siguiente conexión en caso de error
		} else {
			logger.LogsInfo(v.Database)
		}

		// Comprobar la conexión a la base de datos
		if err := dbconn.Ping(); err != nil {
			logger.LogsError(map[string]interface{}{"Error": err})
			continue // Saltar a la siguiente conexión en caso de error
		}

		// Agregar la conexión abierta al slice
		dbConnections = append(dbConnections, dbconn)
	}

	// Retornar el slice de conexiones
	return dbConnections, nil
}
