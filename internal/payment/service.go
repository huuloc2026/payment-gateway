package payment

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Service interface {
	Enqueue(ctx context.Context, p *Payment) error
}

type service struct {
	nc      *nats.Conn
	subject string
	repo    Repository
}

func NewService(nc *nats.Conn, subject string, repo Repository) Service {
	return &service{nc: nc, subject: subject, repo: repo}
}

func (s *service) Enqueue(ctx context.Context, p *Payment) error {
	// persist pending
	if err := s.repo.Create(ctx, p); err != nil {
		return err
	}
	// publish to NATS
	b, _ := json.Marshal(p)
	return s.nc.Publish(s.subject, b)
}
