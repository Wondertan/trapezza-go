//go:generate rm -rf schema
//go:generate go run github.com/99designs/gqlgen

package resolver

import (
	"github.com/Wondertan/trapezza-go/resolver/schema"
	"github.com/Wondertan/trapezza-go/trapezza"
)

type Resolver struct {
	trapezza *trapezza.Manager
}

func NewTrapezzaResolver(trapezza *trapezza.Manager) schema.ResolverRoot {
	return &Resolver{trapezza: trapezza}
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
