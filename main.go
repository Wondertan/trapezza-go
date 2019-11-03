package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"

	"github.com/Wondertan/trapezza-go/resolver"
	"github.com/Wondertan/trapezza-go/restaurant"
	"github.com/Wondertan/trapezza-go/trapezza"
	"github.com/Wondertan/trapezza-go/utils"
)

const (
	defaultPort = "8080"
	endpoint    = "/query"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	signal, ctx := utils.SetupInterruptHandler(context.Background())
	defer signal.Close()

	trapezza := trapezza.NewManager(ctx)
	restaurant := restaurant.NewManager(ctx, trapezza)

	http.Handle("/", handler.Playground("Trapezza playground", endpoint))
	http.Handle(endpoint, resolver.Handler(trapezza, restaurant))

	log.Println("Trapezza-go server start successfully.")
	log.Println("Listening on port: ", defaultPort)
	log.Println("GraphQL enpoint: ", endpoint)

	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
