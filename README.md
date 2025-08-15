# Payment Gateway Simulator (GoFiber + sqlx)

A small, production-like payment gateway to demo at VNPay interviews.

## Stack
- Go 1.22, GoFiber
- sqlx + PostgreSQL
- Redis
- NATS (pub/sub)
- Prometheus + Grafana
- Docker Compose

## Architecture
```mermaid
flowchart LR
    subgraph Client
      M[Merchant]
    end
    subgraph API[API Service]
      R1[POST /payments]
Validate Signature
Save pending
Publish NATS
      R2[GET /payments/:id]
Check Redis -> DB
    end

    subgraph Queue
      N[NATS: payments.created]
    end

    subgraph Processor[Processor Service]
      P[Consume -> Simulate bank
Update DB + Redis]
    end

    DB[(PostgreSQL)]
    C[(Redis)]
    PM[Prometheus]
    GF[Grafana]

    M --> R1 --> N --> P --> DB
    P --> C
    R2 --> C
    R2 --> DB
    API-->PM
    Processor-->PM
    PM-->GF
```

## Run
```bash
cp .env.example .env   # optional, compose sets sane defaults
docker compose up --build -d
```

## Demo
1. Create a signature: `message = orderID|amount` (e.g., `ORD123|100.00`) HMAC-SHA256 with `HMAC_SECRET`.
2. Create payment:
```bash
curl -X POST http://localhost:3000/payments   -H 'Content-Type: application/json'   -d '{"order_id":"ORD123","amount":100.00,"signature":"<hex-signature>"}'
```
3. Poll status:
```bash
curl http://localhost:3000/payments/<id>
```
4. Metrics: Prometheus at http://localhost:9090, Grafana at http://localhost:3001 (default admin/admin).

## Notes
- 80% success, 20% fail (random) to simulate reality.
- Redis caches the latest status for 10 minutes.
- SQL schema is initialized via `deployments/postgres/init.sql`.
```

