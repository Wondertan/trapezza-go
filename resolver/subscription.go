package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/trapezza"
)

type subscription Resolver

func (s *subscription) Updates(ctx context.Context, id string) (<-chan *trapezza.Update, error) {
	ses, err := s.trapezza.SessionByID(id)
	if err != nil {
		return nil, err
	}

	return ses.SubscribeUpdates(ctx)
}
