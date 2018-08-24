package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

func (a *App) Initialize(user, password, dbName string) {
	connectionSpring := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	var err error
	a.DB, err = sql.Open("postgres", connectionSpring)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

}

func (a *App) Run(addr string) {}

