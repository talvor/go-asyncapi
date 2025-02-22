package main

import (
	"log"

	"github.com/talvor/asyncapi/apiserver"
	"github.com/talvor/asyncapi/config"
	"github.com/talvor/asyncapi/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf := config.GetConfig()

	db, err := store.NewPostgresDB(conf)
	if err != nil {
		return err
	}

	dataStore := store.New(db)
	srv := apiserver.New(conf, dataStore)

	err = srv.Start()
	if err != nil {
		return err
	}

	return nil
}
