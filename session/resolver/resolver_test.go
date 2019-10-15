package resolver

import (
	"github.com/99designs/gqlgen/client"
)

func initSession(c *client.Client) string {
	var resp struct {
		NewSession string
	}

	err := c.Post(
		`
			 mutation($waiter: String!, $table: String!) {
				newSession(waiter: $waiter, table: $table) 
			 }
		`,
		&resp,
		client.Var("waiter", "test"),
		client.Var("table", "test"),
	)
	if err != nil {
		panic(err)
	}

	return resp.NewSession
}
