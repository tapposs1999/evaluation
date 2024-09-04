package main

import (
    // "context"
    // "fmt"
    // "time"
    "log"
    "net/http"
    router "evaluation/my-go-project/router"
)

func main() {
    r := router.SetupRouter()
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
