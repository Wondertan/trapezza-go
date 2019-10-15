package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"

	"github.com/Wondertan/trapezza-go/session"
	"github.com/Wondertan/trapezza-go/session/resolver"
)

const defaultPort = "8080"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx := context.Background()
	man := session.NewManager(ctx)

	http.Handle("/", handler.Playground("GraphQL playground", "/session"))
	http.Handle("/session", resolver.Handler(man))

	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
