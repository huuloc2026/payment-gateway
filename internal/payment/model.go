package payment

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

type Payment struct {
	ID        string    `db:"id" json:"id"`
	OrderID   string    `db:"order_id" json:"order_id"`
	Amount    float64   `db:"amount" json:"amount"`
	Status    Status    `db:"status" json:"status"`
	Signature string    `db:"signature" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
