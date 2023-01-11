package CorsConfig

import (
	"encoding/json"
	"os"

	"github.com/gin-contrib/cors"
)

func ReadConfig(file string) cors.Config {
	var corsConfig cors.Config
	jsonBlob, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonBlob, &corsConfig)
	if err != nil {
		panic(err)
	}

	return corsConfig
}
