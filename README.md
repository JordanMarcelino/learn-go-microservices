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
graph TD
  A[Order Created] --> B[Schedule Reminders]
  A --> C[Start 24h Countdown]
  B --> D[4h Reminder]
  B --> E[12h Reminder]
  B --> F[22h Reminder]
  C --> G[Expiration Check]

  D --> H[Send Email via Mail Service]
  E --> H
  F --> H
  G --> I{Payment Status?}

  I -->|Paid| J[Mark Order Completed]
  I -->|Unpaid| K[Cancel Order]
  K --> L[Restock Inventory]
  K --> M[Send Cancellation Email]
```

### 1. Order Creation with Delayed Messages

```mermaid
sequenceDiagram
  participant Client
  participant Gateway
  participant OrderService
  participant PostgreSQL
  participant RabbitMQ
  participant Kafka

  Client->>Gateway: POST /orders
  Gateway->>OrderService: Forward request
  OrderService->>Redis: SETNX lock:request_id
  Redis-->>OrderService: Lock acquired
  OrderService->>PostgreSQL: Create order (status=PENDING)
  OrderService->>RabbitMQ: Publish delayed messages
  Note over RabbitMQ: 4h, 12h, 22h delays
  OrderService->>Kafka: order.created
  OrderService->>Redis: DEL lock:request_id
  Kafka->>ProductService: Reserve stock
```

### 2. Payment Reminder Execution Flow

```mermaid
sequenceDiagram
  participant RabbitMQ
  participant OrderService
  participant PostgreSQL
  participant MailService

  RabbitMQ->>OrderService: Reminder event
  OrderService->>PostgreSQL: Get order status
  alt Status = PENDING
    OrderService->>MailService: Send payment reminder
    MailService->>User: Email reminder
  else Status = PAID/CANCELLED
    OrderService->>RabbitMQ: Acknowledge message
  end
```

### 3. Order Expiration Handling

```mermaid
sequenceDiagram
  participant RabbitMQ
  participant OrderService
  participant PostgreSQL
  participant Kafka
  participant ProductService

  RabbitMQ->>OrderService: Expiration event (24h)
  OrderService->>PostgreSQL: Check payment status
  alt Unpaid
    OrderService->>PostgreSQL: Update status=CANCELLED
    OrderService->>Kafka: order.cancelled
    Kafka->>ProductService: Restock items
    OrderService->>MailService: Notify cancellation
  else Paid
    OrderService->>PostgreSQL: Update status=COMPLETED
  end
```

## 🧠 Key Design Decisions

1. **Delayed Message Implementation**

    ```mermaid
    graph LR
      A[Order Service] -->|Publish with delay| B[RabbitMQ]
      B -->|x-delayed-message plugin| C[Delayed Exchange]
      C -->|TTL 4h| D[Reminder Queue]
      C -->|TTL 12h| E[Reminder Queue]
      C -->|TTL 22h| F[Reminder Queue]
      C -->|TTL 24h| G[Expiration Queue]
    ```

2. **State Transition Diagram**

    ```mermaid
    stateDiagram-v2
      [*] --> PENDING
      PENDING --> PAID: Payment received
      PENDING --> CANCELLED: 24h timeout
      PAID --> COMPLETED: Order fulfilled
      CANCELLED --> [*]
      COMPLETED --> [*]
    ```

**Key Features**:

-   Redis distributed locking for idempotency
-   Event-driven order processing
-   Delayed message handling with RabbitMQ
-   State transition management with PostgreSQL

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
