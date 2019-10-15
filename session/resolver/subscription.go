package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/session"
)

type subscription Resolver

func (s *subscription) SessionEvent(ctx context.Context, session session.ID) (<-chan session.Event, error) {
	ses, err := s.Manager.Session(session)
	if err != nil {
		return nil, err
	}
	
	return ses.SubscribeEvents()
}
