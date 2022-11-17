package main

import (
	"fmt"
	"librarymanagement/controller"
	"librarymanagement/database"
	"log"
	"net/http"
	//"net/http/httptest"
	//"testing"

	"github.com/gorilla/mux"
	//"github.com/stretchr/testify/assert"
)

var port = 8090

func main() {
	database.Connect()
	database.Migrate()

	router := mux.NewRouter().StrictSlash(true)
	controller.RegisterRoutes(router)
	//fmt.Sprintf("Starting Server on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}


