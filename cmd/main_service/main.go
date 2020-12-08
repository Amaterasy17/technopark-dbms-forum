package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)


func main() {
	router := mux.NewRouter()
	addr := ":8080"
	err := http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Println("error of starting server")
	}
}