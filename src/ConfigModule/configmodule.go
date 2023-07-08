package ConfigModule

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigModule struct {
	Port            string
	GoogleMapAPIKey string
	HEREAPIKey      []string
}

func (c ConfigModule) LoadConfig() ConfigModule {

	var config ConfigModule = ConfigModule{Port: "", GoogleMapAPIKey: "", HEREAPIKey: make([]string, 0)}
	path, err := os.Getwd()
	file, err := os.ReadFile(path + "/config/config.json")
	if err != nil {
		log.Fatal("Error loading file: ", err)
	}
	err = json.Unmarshal(file, &config)

	if err != nil {
		log.Fatal("Error parsing file: ", err)
	}
	return config

}
