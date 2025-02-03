# 🚀 Go Microservices E-Commerce Platform

Welcome to my journey of building a scalable microservices architecture using Golang! This project will implement core microservices patterns with modern technologies.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23.4-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-24.0+-blue.svg)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🌟 Key Features

-   **Event-Driven Architecture** with Kafka & RabbitMQ
-   **Database Sharding** with PostgreSQL Cluster
-   **Observability** with OpenTelemetry & Grafana Stack
-   **Production-Grade** Infrastructure Setup

## 🏗 Architecture Overview

```mermaid
graph TD
  A[Client] --> B[API Gateway]
  subgraph Core Services
    B --> C[Auth Service]
    B --> D[Product Service]
    B --> E[Order Service]
    B --> F[Mail Service]
  end

  subgraph Data Layer
    C --> G[PostgreSQL Cluster]
    D --> H[PostgreSQL Cluster]
    E --> I[PostgreSQL Cluster]
    E --> J[Redis Cluster]
  end

  subgraph Event Bus
    D --> K[Kafka Cluster]
    E --> K
    C --> L[RabbitMQ Cluster]
    F --> L
  end
```

## 🛠 Service Breakdown

### 1. 🔐 Auth Service

**Responsibilities**: User identity management and authentication

```mermaid
sequenceDiagram
  participant Client
  participant Gateway
  participant AuthService
  participant PostgreSQL
  participant RabbitMQ

  Client->>Gateway: POST /register
  Gateway->>AuthService: Forward request
  AuthService->>PostgreSQL: Create user (unverified)
  AuthService->>RabbitMQ: Publish verification event
  RabbitMQ->>MailService: Consume event
  MailService->>User: Send verification email
```

**Key Features**:

-   JWT-based authentication
-   Email verification with 10-minute expiry
-   Anti-spam protection (1-minute cooldown)
-   CQRS pattern with master-slave replication
-   Horizontal scaling with pgpool-II

---

### 2. 📦 Product Service

**Responsibilities**: Product lifecycle management

```mermaid
graph TD
  A[Create Product] --> B[Kafka: product-created]
  C[Update Product] --> D[Kafka: product-updated]
  E[Order Created Event] --> F[Reduce Stock]
  B --> G[Order Service]
  D --> G
```

**Key Features**:

-   Real-time inventory synchronization
-   Event sourcing for product changes
-   CQRS pattern with master-slave replication

---

### 3. 🛒 Order Service

**Responsibilities**: Order processing and fulfillment

```mermaid
sequenceDiagram
  participant Client
  participant Gateway
  participant OrderService
  participant Redis
  participant Kafka

  Client->>Gateway: POST /orders
  Gateway->>OrderService: Forward request
  OrderService->>Redis: Acquire lock (request ID)
  Redis-->>OrderService: Lock acquired
  OrderService->>Kafka: Publish order-created
  Kafka->>ProductService: Update inventory
  OrderService->>PostgreSQL: Commit order
  OrderService->>Redis: Release lock
```

**Key Features**:

-   Redis distributed locking for idempotency
-   Event-driven order processing
-   Circuit breaker pattern for inventory checks

---

### 4. 📧 Mail Service

**Responsibilities**: Asynchronous email processing

**Key Features**:

-   RabbitMQ consumer
-   Template-based email rendering
-   Send retry mechanism with exponential backoff
-   MailHog integration for development

---

### 5. 🌉 API Gateway

**Responsibilities**: Unified API entrypoint

**Key Features**:

-   JWT validation middleware
-   Rate limiting per service
-   Request/Response transformation
-   Prometheus metrics collection

## 🏭 Infrastructure Architecture

```mermaid
graph TD
  subgraph Database Layer
    A[PostgreSQL Cluster] --> B[Auth DB]
    A --> C[Product DB]
    A --> D[Order DB]
    B --> E[1 Master + 2 Replicas]
    C --> F[1 Master + 2 Replicas]
    D --> G[1 Master + 2 Replicas]
  end

  subgraph Message Brokers
    H[RabbitMQ Cluster] --> I[3 Nodes]
    J[Kafka Cluster] --> K[3 Brokers]
  end

  subgraph Caching
    L[Redis Cluster] --> M[3 Masters + 3 Replicas]
  end

  subgraph Observability
    N[Prometheus]
    O[Loki]
    P[Tempo]
    Q[Grafana]
    R[OpenTelemetry]
  end
```

## 🔍 Observability Stack

```mermaid
graph TD
  A[Services] --> B[OpenTelemetry]
  B --> C[Metrics]
  B --> D[Traces]
  B --> E[Logs]
  C --> F[Prometheus]
  D --> G[Tempo]
  E --> H[Loki]
  F --> I[Grafana]
  G --> I
  H --> I
```

## 📈 Deployment Architecture

```mermaid
graph TD
  A[Load Balancer] --> B[API Gateway Cluster]
  B --> C[Auth Service Cluster]
  B --> D[Product Service Cluster]
  B --> E[Order Service Cluster]
  B --> F[Mail Service Cluster]

  C --> G[PostgreSQL Cluster]
  D --> H[PostgreSQL Cluster]
  E --> I[PostgreSQL Cluster]
  E --> J[Redis Cluster]

  C --> K[RabbitMQ Cluster]
  D --> L[Kafka Cluster]
  E --> L
  F --> K
```

**Monitoring Features**:

-   Real-time service metrics
-   Distributed tracing across services
-   Centralized logging with labels
-   Performance dashboards per service

## 🛠️ Technology Stack

**Languages & Frameworks**

-   Go 1.23+
-   Gin Web Framework

**Databases**

-   PostgreSQL 16 with pgpool-II
-   Bitnami Redis Latest Cluster

**Message Brokers**

-   Bitnami Kafka Latest Cluster
-   RabbitMQ 4.0 Cluster

**Infrastructure**

-   Docker Swarm
-   HAProxy for load balancing
-   MailHog SMTP server

**Observability**

-   Prometheus
-   Grafana
-   Loki
-   Tempo
-   OpenTelemetry
