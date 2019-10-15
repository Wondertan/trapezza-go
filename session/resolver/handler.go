package resolver

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/websocket"

	"github.com/Wondertan/trapezza-go/session"
	"github.com/Wondertan/trapezza-go/session/resolver/schema"
)

func Handler(manager *session.Manager) http.HandlerFunc {
	return handler.GraphQL(
		schema.NewExecutableSchema(
			schema.Config{
				Resolvers: NewSessionResolver(manager),
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
