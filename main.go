package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go-jwt/dbutils"
	"go-jwt/handler"
	"go-jwt/migrations"
	"log"
	"net/http"
)

var serverAddr = "localhost:4000"
var dbPath = "./go-jwt.db"

func main() {
	db, dbClose := dbutils.ConnectSqlite(dbPath)
	defer dbClose()

	migrations.CreateTables(db)

	appCtx := &handler.AppContext{DB: db, Validator: validator.New()}
	router := mux.NewRouter()

	handler.RegisterRoutes(router, appCtx)

	log.Printf("starting the server at %s \n", serverAddr)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Panicf("failed to start server: %s \n", err.Error())
	}
}
