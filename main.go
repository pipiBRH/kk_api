package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pipiBRH/kk_api/httpserver"

	"github.com/pipiBRH/kk_api/config"
	"github.com/pipiBRH/kk_api/database"
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

	// esDAL := dal.NewElasticsearchDAL(database.EsClient)
	// err = esDAL.SearchMaxAndMinFreeSpaceYouBikeInfoByTimeRange(
	// 	config.Config.Elasticsearch.Index["youbie"],
	// 	"2019-06-29 12:00:00",
	// 	"2019-06-29 23:00:00")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(data)

	httpserver.NewServer(config.Config.App.Host, config.Config.App.Port, gin.Default()).InitServer()
}
