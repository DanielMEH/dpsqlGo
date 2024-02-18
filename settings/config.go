package settings

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const (
	AND                   string = " AND "
	OR                    string = " OR "
	NOT                   string = " NOT "
	EQUAL                 string = " = "
	NOT_EQUAL             string = " != "
	LIKE                  string = " LIKE "
	LESS_THAN             string = " < "
	GREATER_THAN          string = " > "
	LESS_THAN_OR_EQUAL    string = " <= "
	GREATER_THAN_OR_EQUAL string = " >= "
	BETWEEN               string = " BETWEEN "
	WHERE                 string = " WHERE "
	LIMIT                string = "LIMIT"
)

type Config struct {
	Log struct {
		FileName string `json:"FileName"`
	} `json:"log"`

	ConfigDb struct {
		Port      string `json:"Port"`
		Protocolo string `json:"Protocolo"`
		Host      string `json:"Host"`
	} `json:"configDb"`

	DbListConn []struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"databases"`
		Schema   string `json:"schema"`
		Port     string `json:"port"`
	} `json:"Db_list_conn"`
}

func LoadConfig() (*Config, error) {
	// Obtener la ruta del directorio del paquete utils
	packageDir, err := filepath.Abs(".")
	if err != nil {

		return nil, err
	}
	// Construir la ruta al archivo JSON y moverse un directorio hacia arriba
	jsonFilePath := filepath.Join(packageDir, "data", "svc_settings.json")

	// Leer y deserializar la configuración desde el archivo JSON
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
