package main

import (
	"context"
	"flag"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"net/http"
	"ozon-broker/internal/dal/db"
	"ozon-broker/internal/handler"
	"sync"

	"log"
	"time"

	"github.com/diphantxm/ozon-api-client/ozon"
)

const SCOPE = drive.DriveScope

func main() {

	log.Printf("Starting")

	configFile := flag.String("c", "config.yaml", "specify path to a config.yaml")
	flag.Parse()
	cfg, err := configure(*configFile)
	if err != nil {
		log.Fatalf("couldn't read config: %s", err)

		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	psql, err := db.InitDB(cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password, cfg.Db.Dbname)
	if err != nil {
		log.Fatalf("Cannot init DB: %s", err)
	}
	defer psql.Close()
	client := ozon.NewClient(http.DefaultClient, cfg.Ozon.Id, cfg.Ozon.Key)
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(cfg.Google.File), option.WithScopes(SCOPE))
	if err != nil {
		log.Fatalf("Warning: Unable to create drive Client %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go handler.OrderComposer(ctx, client, psql, &wg, srv)
	go handler.ProductComposer(ctx, client, psql, &wg, srv)
	wg.Wait()

	log.Printf("Finished")

}
