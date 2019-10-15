//go:generate rm -rf schema
//go:generate go run github.com/99designs/gqlgen

package resolver

import (
	"github.com/Wondertan/trapezza-go/session"
	"github.com/Wondertan/trapezza-go/session/resolver/schema")

type Resolver struct {
	*session.Manager
}

func NewSessionResolver(man *session.Manager) schema.ResolverRoot {
	return &Resolver{Manager: man}
}

func (r *Resolver) Mutation() schema.MutationResolver {
	return (*mutation)(r)
}

func (r *Resolver) Query() schema.QueryResolver {
	return (*query)(r)
}

func (r *Resolver) Subscription() schema.SubscriptionResolver {
	return (*subscription)(r)
}

