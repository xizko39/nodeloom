package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("Starting NodeLoom backend server...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
