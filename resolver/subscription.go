package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/restaurant"
	"github.com/Wondertan/trapezza-go/trapezza"
)

type subscription Resolver

func (s *subscription) TrapezzaSessionUpdates(ctx context.Context, id string) (<-chan *trapezza.Update, error) {
	ses, err := s.trapezza.Session(id)
	if err != nil {
		return nil, err
	}

	return ses.SubscribeUpdates(ctx)
}

func (s *subscription) RestaurantEvents(ctx context.Context, id string) (<-chan restaurant.Event, error) {
	return s.restaurant.SubscribeEvents(ctx, id)
}
