package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"

	"github.com/Wondertan/trapezza-go/resolver"
	"github.com/Wondertan/trapezza-go/trapezza"
)

const defaultPort = "8080"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx := context.Background()
	man := trapezza.NewManager(ctx)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", resolver.Handler(man))

	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}

// func newClient() *dgo.Dgraph {
// 	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	return dgo.NewDgraphClient(
// 		api.NewDgraphClient(d),
// 	)
// }
