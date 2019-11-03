package resolver

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/websocket"

	"github.com/Wondertan/trapezza-go/resolver/schema"
	"github.com/Wondertan/trapezza-go/restaurant"
	"github.com/Wondertan/trapezza-go/trapezza"
)

func Handler(trapezza *trapezza.Manager, restaurant *restaurant.Manager) http.HandlerFunc {
	return handler.GraphQL(
		schema.NewExecutableSchema(
			schema.Config{
				Resolvers: NewTrapezzaResolver(trapezza, restaurant),
			},
		),
		handler.WebsocketUpgrader(
			websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		),
		handler.WebsocketKeepAliveDuration(10*time.Second),
	)
}
