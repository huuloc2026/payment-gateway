
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/yourname/payment-gateway-simulator/internal/cache"
	"github.com/yourname/payment-gateway-simulator/internal/config"
	"github.com/yourname/payment-gateway-simulator/internal/db"
	"github.com/yourname/payment-gateway-simulator/internal/payment"
	"github.com/yourname/payment-gateway-simulator/internal/queue"
)

func main() {
	_ = godotenv.Load()
	cfg := config.New()

	database := db.Connect(cfg.DB.DSN)
	repo := payment.NewRepository(database)
	nc := queue.Connect(cfg.NATS.URL)
	rdb := cache.New(cfg.Redis.Addr, cfg.Redis.DB)

	// metrics
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Processor.MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		log.Info().Msgf("processor metrics on %s", addr)
		_ = http.ListenAndServe(addr, nil)
	}()

	_, err := nc.Subscribe(cfg.NATS.Subject, func(msg *nats.Msg) {
		var p payment.Payment
		if err := json.Unmarshal(msg.Data, &p); err != nil {
			log.Error().Err(err).Msg("unmarshal")
			return
		}
		process(ctxWithTimeout(), &p, repo, rdb)
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	select {}
}

func ctxWithTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func process(ctx context.Context, p *payment.Payment, repo payment.Repository, rdb *redis.Client) {
	// simulate bank processing latency
	time.Sleep(time.Duration(500+rand.Intn(1200)) * time.Millisecond)
	// random outcome (80% success)
	if rand.Intn(100) < 80 {
		_ = repo.UpdateStatus(ctx, p.ID, payment.StatusSuccess)
		_ = rdb.Set(ctx, "payment:"+p.ID, string(payment.StatusSuccess), 10*time.Minute).Err()
	} else {
		_ = repo.UpdateStatus(ctx, p.ID, payment.StatusFailed)
		_ = rdb.Set(ctx, "payment:"+p.ID, string(payment.StatusFailed), 10*time.Minute).Err()
	}
}
