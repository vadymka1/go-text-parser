package main

import (
	"fmt"
	"github.com/go-text-parse/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Starting ...")

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.GetuploadForm).Methods("GET")
	router.HandleFunc("/upload", controllers.GetStatistic).Methods("POST")

	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal(err)
	}

}
