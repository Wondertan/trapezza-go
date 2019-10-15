package resolver

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"

	"github.com/Wondertan/trapezza-go/session"
)

func TestQuery(t *testing.T) {
	c := client.New(Handler(session.NewManager(context.Background())))
	id := initSession(c)

	t.Run("Session", func(t *testing.T) {
		var resp struct {
			Session struct {
				Id     string
				Waiter string
				Table  string
				Orders []struct {
					Client string
					Items  []string
				}
			}
		}

		err := c.Post(
			`
				 query($id: ID!) {
					session(id: $id) {
						id,
						waiter,
						table,
						orders {
							client,
							items
						}
					}
				 }
			`,
			&resp,
			client.Var("id", id),
		)
		if err != nil {
			t.Fatal(err)
		}
	})
}
