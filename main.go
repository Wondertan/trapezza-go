package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"

	"github.com/Wondertan/trapezza-go/schema"
	"github.com/Wondertan/trapezza-go/session"
)

const defaultPort = "8080"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx := context.Background()
	man := session.NewManager(ctx)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(schema.NewExecutableSchema(schema.Config{Resolvers: &schema.Resolver{Manager: man}})))

	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
