package queue

import (
	"log"

	"github.com/nats-io/nats.go"
)

func Connect(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect nats: %v", err)
	}
	return nc
}
