package router

import (
    "github.com/gorilla/mux"
	pricecontroller "evaluation/my-go-project/controller"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/api/get", pricecontroller.Get).Methods("GET")
    r.HandleFunc("/api/insert", pricecontroller.Insert).Methods("POST")
    return r
}
