# eCommerce CQRS/ES in Go

A simple eCommerce site demonstrating Event Sourcing (ES) and Command Query Responsibility Segregation (CQRS).

## Architecture

- **Command Side (Write)**: Validates business logic by reconstructing the aggregate state from the event store and persists new events.
- **Event Store**: A PostgreSQL table that stores all state changes as immutable events.
- **Query Side (Read)**: A separate schema optimized for fast reads, updated by a projectionist that listens to events.
- **Projectionist**: Synchronizes the read models when new events are saved.

## Tech Stack

- **Backend**: Go 1.24
- **Database**: PostgreSQL
- **Orchestration**: Docker Compose

## Getting Started

### 1. Start the services
```bash
docker compose up -d
```

### 2. Manual Schema Setup (For this demo)
Since we don't have a migration tool, you'll need to apply the schemas manually if the app doesn't do it (note: main.go needs to be updated to apply schemas if desired, but for this demo, you can run them via psql or the app can be modified).

*Actually, let's update main.go to apply the schemas for a better experience.*

### API Endpoints

- `POST /products`: Create a product
  ```json
  { "id": "p1", "name": "Laptop", "price": 1200.00 }
  ```
- `POST /inventory/adjust`: Adjust stock
  ```json
  { "product_id": "p1", "quantity": 10 }
  ```
- `POST /orders`: Place an order
  ```json
  { "order_id": "o1", "product_id": "p1", "quantity": 1 }
  ```
- `GET /products`: List all products
- `GET /inventory/history/{product_id}`: View inventory changes
- `GET /orders`: View order history

---

## 🛠️ Project Deep Dive & Technical Guide

### 🎯 Main Purpose
High-performance E-commerce microservice backend in Go utilizing CQRS pattern.

### 🧰 Tecnologias Utilizadas (Tech Stack)
- **Primary Technologies:** Docker / Docker Compose, Go 1.24.3
- **Deployment Strategy:** Container-ready with standard Docker/Docker Compose integrations where applicable.

### 🚀 Como Implantar e Executar (How to Deploy & Run)
#### Go CLI / Service
1. Build the module: `go build`.
2. Execute: `./ecommerce-cqrs-go` (or project executable name).

---
