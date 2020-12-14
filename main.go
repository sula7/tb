package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"

	"tb/api"
	"tb/storage"
)

type config struct {
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PWD,required"`
}

func main() {
	var conf config
	err := env.Parse(&conf)
	if err != nil {
		log.Fatal("unable to get env config:", err)
	}

	store, err := storage.New(fmt.Sprintf(`postgres://%s:%s@localhost:5432/tochka`, conf.DbUser, conf.DbPassword))
	if err != nil {
		log.Fatal("unable to establish DB connection:", err)
	}

	defer store.Close()

	err = store.Migrate()
	if err != nil {
		log.Fatal("unable to apply migrations:", err)
	}

	apiService := api.New(store)
	apiService.Start()
}
