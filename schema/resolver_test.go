package schema

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"

	"github.com/Wondertan/trapezza-go/session"
)

func TestMutationResolver(t *testing.T) {
	ctx := context.Background()
	man := session.NewManager(ctx)
	c := client.New(handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{Manager: man}})))

	var resp struct {
		NewSession string
	}

	err := c.Post(
		`
			 mutation($waiter: ID!, $table: ID!) {
				newSession(waiter: $waiter, table: $table) 
			 }
		`,
		&resp,
		client.Var("waiter", "test"),
		client.Var("table", "test"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if resp.NewSession == "" {
		t.Fatal("zero id")
	}

	id := resp.NewSession

	t.Run("AddClient", func(t *testing.T) {
		var resp struct {
			AddClient bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $client: ID!) {
					addClient(session: $id, client: $client) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("client", "test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.AddClient {
			t.Fatal("response is not true")
		}
	})

	t.Run("AddItem", func(t *testing.T) {
		var resp struct {
			AddItem bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $client: ID!, $item: ID!) {
					addItem(session: $id, client: $client, item: $item) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("client", "test"),
			client.Var("item", "test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.AddItem {
			t.Fatal("response is not true")
		}
	})

	t.Run("SetWaiter", func(t *testing.T) {
		var resp struct {
			SetWaiter bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $waiter: ID!) {
					setWaiter(session: $id, waiter: $waiter) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("waiter", "test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.SetWaiter {
			t.Fatal("response is not true")
		}
	})

	t.Run("SetTable", func(t *testing.T) {
		var resp struct {
			SetTable bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $table: ID!) {
					setTable(session: $id, table: $table) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("table", "test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.SetTable {
			t.Fatal("response is not true")
		}
	})

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

	t.Run("EndSession", func(t *testing.T) {
		var resp struct {
			EndSession bool
		}

		err := c.Post(
			`
				 mutation($id: ID!) {
					endSession(session: $id) 
				 }
			`,
			&resp,
			client.Var("id", id),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.EndSession {
			t.Fatal("response is not true")
		}
	})

	t.Run("SessionError", func(t *testing.T) {
		err := c.Post(
			`
				 mutation($id: ID!) {
					endSession(session: $id) 
				 }
			`,
			nil,
			client.Var("id", id),
		)
		if err == nil || err.Error() != "[{\"message\":\"session: not found\",\"path\":[\"endSession\"]}]" {
			t.Error("wrong or zero error")
		}
	})
}
