package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic"
)

type ElasticsearchConnection struct {
	Ctx    context.Context
	Client *elastic.Client
}

var EsClient *ElasticsearchConnection

func InitElasticsearchConnection(host string, port int) error {
	ctx := context.Background()
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	addr := fmt.Sprintf("http://%s:%d", host, port)

	client, err := elastic.NewClient(
		elastic.SetErrorLog(errorlog),
		elastic.SetURL(addr),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		return err
	}

	_, _, err = client.Ping(addr).Do(ctx)
	if err != nil {
		return err
	}

	EsClient = &ElasticsearchConnection{
		Ctx:    ctx,
		Client: client,
	}

	return nil
}
