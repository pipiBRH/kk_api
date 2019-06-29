package main

import (
	"github.com/pipiBRH/kk_database"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pipiBRH/kk_api/httpserver"

	"github.com/pipiBRH/kk_api/config"
)

func main() {
	err := config.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitElasticsearchConnection(
		config.Config.Elasticsearch.Host,
		config.Config.Elasticsearch.Port,
	)
	if err != nil {
		log.Fatal(err)
	}

	httpserver.NewServer(config.Config.App.Host, config.Config.App.Port, gin.Default()).InitServer()
}
