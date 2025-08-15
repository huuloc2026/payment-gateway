package payment

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, p *Payment) error
	UpdateStatus(ctx context.Context, id string, status Status) error
	GetByID(ctx context.Context, id string) (*Payment, error)
}

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) Repository { return &repo{db} }

func (r *repo) Create(ctx context.Context, p *Payment) error {
	q := `INSERT INTO payments (id, order_id, amount, status, signature, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,NOW(),NOW())`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.OrderID, p.Amount, p.Status, p.Signature)
	return err
}

func (r *repo) UpdateStatus(ctx context.Context, id string, status Status) error {
	q := `UPDATE payments SET status=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id, status)
	return err
}

func (r *repo) GetByID(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	q := `SELECT id, order_id, amount, status, signature, created_at, updated_at FROM payments WHERE id=$1`
	err := r.db.GetContext(ctx, &p, q, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
