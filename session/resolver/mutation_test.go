package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"

	"github.com/Wondertan/trapezza-go/session"
)

func TestMutationSubscription(t *testing.T) {
	c := client.New(Handler(session.NewManager(context.Background())))
	id := initSession(c)

	sub := c.Websocket(
		`
		subscription($id: ID!) {
			sessionEvent(id: $id) {
				... on ClientEvent {
					client
				}
				... on ItemEvent {
					client
					item
				}
				... on WaiterEvent {
					waiter
				}
				... on TableEvent {
					table
				}
			}
		}
		`,
		client.Var("id", id),
	)
	defer sub.Close()

	time.Sleep(10 * time.Millisecond) // give time for subscription to init

	t.Run("AddClient", func(t *testing.T) {
		in := "test"

		var resp struct {
			AddClient bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $client: String!) {
					addClient(session: $id, client: $client) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("client", in),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.AddClient {
			t.Fatal("response is not true")
		}

		var event struct {
			SessionEvent struct{
				Client string
			}
		}

		err = sub.Next(&event)
		if err != nil {
			t.Fatal(err)
		}

		if event.SessionEvent.Client != in {
			t.Fatal("clients are not equal")
		}
	})

	t.Run("AddItem", func(t *testing.T) {
		in := "test"
		item := "test"

		var resp struct {
			AddItem bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $client: String!, $item: String!) {
					addItem(session: $id, client: $client, item: $item) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("client", in),
			client.Var("item", item),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.AddItem {
			t.Fatal("response is not true")
		}

		var event struct {
			SessionEvent struct{
				Client string
				Item   string
			}
		}

		err = sub.Next(&event)
		if err != nil {
			t.Fatal(err)
		}

		if event.SessionEvent.Client != in {
			t.Fatal("clients are not equal")
		}

		if event.SessionEvent.Item != item {
			t.Fatal("items are not equal")
		}
	})

	t.Run("SetWaiter", func(t *testing.T) {
		waiter := "test"

		var resp struct {
			SetWaiter bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $waiter: String!) {
					setWaiter(session: $id, waiter: $waiter) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("waiter", waiter),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.SetWaiter {
			t.Fatal("response is not true")
		}

		var event struct {
			SessionEvent struct{
				Waiter string
			}
		}

		err = sub.Next(&event)
		if err != nil {
			t.Fatal(err)
		}

		if event.SessionEvent.Waiter != waiter {
			t.Fatal("waiters are not equal")
		}
	})

	t.Run("SetTable", func(t *testing.T) {
		table := "test"

		var resp struct {
			SetTable bool
		}

		err := c.Post(
			`
				 mutation($id: ID!, $table: String!) {
					setTable(session: $id, table: $table) 
				 }
			`,
			&resp,
			client.Var("id", id),
			client.Var("table", table),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !resp.SetTable {
			t.Fatal("response is not true")
		}

		var event struct {
			SessionEvent struct{
				Table string
			}
		}

		err = sub.Next(&event)
		if err != nil {
			t.Fatal(err)
		}

		if event.SessionEvent.Table != table {
			t.Fatal("tables are not equal")
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
