package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"encoding/json"
)

type App struct {
	Router     *mux.Router
	DB         *sql.DB
}

func (a *App) Initialize(user, password, dbName string) {
	connectionSpring := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	var err error
	a.DB, err = sql.Open("postgres", connectionSpring)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
}

func (a *App) getProduct(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid Product Id")
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(writer, http.StatusNotFound, "Product not found")
		default:
			respondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return	}

	respondWithJson(writer, http.StatusOK, p)
}

func respondWithError(writer http.ResponseWriter, code int, message string  ) {
	respondWithJson(writer, code, map[string]string{"error": message})
}

func respondWithJson(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}
