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
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
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

func (a *App) getProducts(writer http.ResponseWriter, req *http.Request) {
	count, _ := strconv.Atoi(req.FormValue("count"))
	start, _ := strconv.Atoi(req.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(writer, http.StatusOK, products)
}


func (a *App) createProduct(writer http.ResponseWriter, req *http.Request) {
	var product product
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&product); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer req.Body.Close()

	if err := product.createProduct(a.DB) ;err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(writer, http.StatusCreated, product)
}

func (a *App) updateProduct(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid Product Id")
	}

	var product product
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&product); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer req.Body.Close()
	product.ID = id


	if err := product.updateProduct(a.DB) ;err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(writer, http.StatusOK, product)
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
