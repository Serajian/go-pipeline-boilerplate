# Go Pipeline Boilerplate

A production-grade **Go boilerplate** for building scalable applications based on the **Pipeline pattern**.  
It provides a clean architecture for composing stages, running them either in parallel, sequential (short-circuit), or barrier mode, and integrates with HTTP and Kafka out of the box.

---

## ✨ Features
- **Pipeline architecture** with three execution strategies:
    - `ChainPipeline` → Parallel stages with concurrent processing.
    - `ShortCircuitPipeline` → Sequential stages that stop immediately on error.
    - `BarrierPipeline` → Parallel processing that waits for all results before continuing.
- **Stage abstraction** (`Stage[T]` and `StageFunc[T]`) for reusability.
- **Clean Dependency Injection (DI)** containers for stages and pipelines.
- **Infrastructure adapters** for:
    - HTTP server (Gin-based).
    - Kafka producer/consumer (Sarama-based).
- **Graceful shutdown** with context propagation.
- **Structured logging** with zero-log.
- **Configuration** via YAML (`config/config.yaml`).
- **Deployment-ready** with Docker, docker-compose, Prometheus, Grafana, Redis, and Postgres.

---

## 📂 Project Structure

```bash
.
├── bootstrap/            # Application bootstrap (initialization & lifecycle)
├── config/               # Config loader & constants
├── deployment/           # Docker, compose, monitoring configs
├── infrastructure/       # Adapters: httpserver, message_queue, registry
├── internal/
│   ├── di/               # Dependency injection containers
│   ├── model/            # Domain models
│   ├── pipelines/        # Pipeline runners (parallel, short-circuit, barrier)
│   ├── ports/            # Interfaces (contracts)
│   ├── presentation/     # HTTP handlers, MQ consumers
│   └── stages/           # Concrete stage implementations
├── pkg/                  # Shared utilities: logger, apperror, generator
├── main.go               # Entry point
└── Makefile              # Lint, build, run helpers
---

