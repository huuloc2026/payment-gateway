
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yourname/payment-gateway-simulator/internal/cache"
	"github.com/yourname/payment-gateway-simulator/internal/config"
	"github.com/yourname/payment-gateway-simulator/internal/db"
	"github.com/yourname/payment-gateway-simulator/internal/payment"
	"github.com/yourname/payment-gateway-simulator/internal/queue"
	"github.com/yourname/payment-gateway-simulator/pkg/utils"
)

func main() {
	_ = godotenv.Load()
	cfg := config.New()

	zerolog.TimeFieldFormat = time.RFC3339

	// infra
	database := db.Connect(cfg.DB.DSN)
	rdb := cache.New(cfg.Redis.Addr, cfg.Redis.DB)
	nc := queue.Connect(cfg.NATS.URL)
	repo := payment.NewRepository(database)
	svc := payment.NewService(nc, cfg.NATS.Subject, repo)

	// metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf(":%s", cfg.API.MetricsPort)
		log.Info().Msgf("metrics listening on %s", addr)
		_ = http.ListenAndServe(addr, nil)
	}()

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		ctx := context.Background()
		if err := cache.Ping(ctx, rdb); err != nil {
			return c.Status(500).JSON(fiber.Map{"redis": "down"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	type CreatePaymentReq struct {
		OrderID   string  `json:"order_id"`
		Amount    float64 `json:"amount"`
		Signature string  `json:"signature"`
	}

	app.Post("/payments", func(c *fiber.Ctx) error {
		var req CreatePaymentReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		if req.OrderID == "" || req.Amount <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid order_id/amount")
		}

		// verify signature (message=order_id|amount with 2 decimal places)
		msg := fmt.Sprintf("%s|%.2f", req.OrderID, req.Amount)
		expected := utils.HMACSHA256Hex(msg, cfg.Security.HMACSecret)
		if expected != req.Signature {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid signature")
		}

		id := uuid.New().String()
		p := &payment.Payment{
			ID:        id,
			OrderID:   req.OrderID,
			Amount:    req.Amount,
			Status:    payment.StatusPending,
			Signature: req.Signature,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := svc.Enqueue(c.Context(), p); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(201).JSON(fiber.Map{
			"id":      id,
			"message": "Payment created, processing...",
		})
	})

	app.Get("/payments/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		// try cache first
		if st, err := rdb.Get(c.Context(), "payment:"+id).Result(); err == nil {
			return c.JSON(fiber.Map{"id": id, "status": st})
		}
		p, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(p)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = cfg.API.Port
	}
	log.Info().Msgf("API listening on :%s", port)
	log.Fatal().Err(app.Listen(":" + port)).Send()
}
