package resolver

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/websocket"

	"github.com/Wondertan/trapezza-go/resolver/schema"
	"github.com/Wondertan/trapezza-go/trapezza"
)

func Handler(manager *trapezza.Manager) http.HandlerFunc {
	return handler.GraphQL(
		schema.NewExecutableSchema(
			schema.Config{
				Resolvers: NewTrapezzaResolver(manager),
			},
		),
		handler.WebsocketUpgrader(
			websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		),
	)
}
