![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)
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
```
---

## 🗺️ Architecture at a Glance
````mermaid
flowchart TD
  A[main.go]
  B[bootstrap: Initialize / Start / Stop]
  C[DI Containers]
  D[MQ Registry]
    D1[Consumer Handler<br/>presentation/mq]
    D2[Run Pipeline<br/> Chain or ShortCircuit or Barrier]
    D3[Stage<br/> Validation]
    D4[Stage<br/> ...]
    D5[Stage Produce]
  E[HTTPServer Registry]
    E1[HTTP Layer<br/>Gin Adapter  presentation/http]
    E2[Request Handlers + Middleware]
    E3[Run Pipeline<br/> Chain or ShortCircuit or Barrier]
    E4[Stage<br/> Auth]
    E5[Stage<br/> ...]

%% Main
A --> B

%% Dependencies wired by bootstrap
B --> C
B --> D
B --> E

%% Messaging
D --> D1
D1 --> D2
D2 --> D3
D3 --> D4
D4 --> D5

%% HTTP
E --> E1
E1 --> E2
E2 --> E3
E3 --> E4
E4 --> E5

%% DI
C --> C1[Pipeline Runners<br/> Chain<br/> ShortCircuit<br/>  Barrier]
C --> C2[Stages<br/> Validation<br/> ... ]

````
### 🧪 Pipeline Execution Modes
#### 1) Chain (Concurrent + Streamed)
````mermaid
sequenceDiagram
  autonumber
  participant H as HTTP/MQ Handler
  participant P as ChainPipeline
  participant S1 as Stage A
  participant S2 as Stage B
  participant S3 as Stage C

  H->>P: Run(ctx, in<-chan T)
  P->>S1: start goroutine + select on ctx/in
  S1-->>P: out1, err1<-chan
  P->>S2: start goroutine (consume out1)
  S2-->>P: out2, err2<-chan
  P->>S3: start goroutine (consume out2)
  S3-->>P: out3, err3<-chan
  P-->>H: merge(out3), mergeErrors(err1,err2,err3)

````
````mermaid
flowchart LR
  IN[<-chan T] --> S1[Stage A]
  S1 --> S2[Stage B]
  S2 --> S3[Stage C]
  S3 --> OUT[<-chan T]

  ERR1[err chan] -.-> OUTERR
  ERR2[err chan] -.-> OUTERR
  ERR3[err chan] -.-> OUTERR

  OUTERR([merged errors])

````
هر استیج در گوروتین خودش با کانال‌ها کار می‌کند؛ ارورها در یک کانال merge می‌شوند.

مناسب برای throughput بالا با backpressure.

#### 2) Short-Circuit (Sequential + Early Exit)
````mermaid
flowchart TD
  IN[Input T] --> S1[StageFn A ctx, T -> T, err]
  S1 -->|ok| S2[StageFn B]
  S1 -->|err| STOP([return err])
  S2 -->|ok| S3[StageFn C]
  S2 -->|err| STOP
  S3 -->|ok| OUT[Final T]
  S3 -->|err| STOP
````
stages به‌ترتیب و blocking اجرا می‌شوند؛ با اولین خطا متوقف می‌شویم.

#### 3) Barrier (Parallel + Wait-All)
````mermaid
flowchart LR
  IN[Input T] --> BARR(Spawn N Stages)
  BARR --> A[Stage A]
  BARR --> B[Stage B]
  BARR --> C[Stage C]
  A --> J[Join/Barrier]
  B --> J
  C --> J
  J --> DEC{All OK?}
  DEC -->|Yes| OUT[Reduce/Merge Results]
  DEC -->|No| ERR[Aggregate Errors]

````
````mermaid
flowchart TD
  IN[<-chan T] --> FANOUT{spawn goroutines}

  FANOUT --> A[Stage A]
  FANOUT --> B[Stage B]
  FANOUT --> C[Stage C]

  A --> JOIN
  B --> JOIN
  C --> JOIN

  JOIN --> DEC{all results ready}
  DEC -->|ok| OUT[finalOut <-chan T]
  DEC -->|err| OUTERR[mergedErr <-chan error]

````
همهٔ stages روی یک ورودی به‌صورت موازی اجرا می‌شوند؛ بعد از جمع شدن نتایج/خطاها ادامه می‌دهیم.
---
