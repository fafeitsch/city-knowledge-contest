package main

import (
	"github.com/fafeitsch/city-knowledge-contest/backend/webapi"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	source := rand.NewSource(time.Now().UnixMilli())
	random := rand.New(source)
	handler := webapi.HandleFunc(webapi.Options{Random: random})
	log.Fatal(http.ListenAndServe(":23123", handler))
}
