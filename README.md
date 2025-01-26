# 🚀 Learn Go Microservices Project

Welcome to my journey of building a scalable microservices architecture using Golang! This project will implement core microservices patterns with modern technologies.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23.4-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-24.0+-blue.svg)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📋 Project Overview

### What We're Building

A simple e-commerce platform demonstrating microservices fundamentals:

-   **Core Services**
    -   🛂 API Gateway
    -   👥 Auth Service
    -   📦 Product Service
    -   🛒 Order Service
    -   📊 Analytics Service

### Key Architecture Components

```mermaid
graph TD
  A[Client] --> B[API Gateway]
  B --> C[Auth Service]
  B --> D[Product Service]
  B --> E[Order Service]
  B --> L[Mail Service]
  C --> F[PostgreSQL]
  D --> G[PostgreSQL]
  D --> J[Kafka]
  E --> H[Redis]
  E --> I[Kafka]
  E --> K[RabbitMQ]
  L --> M[RabbitMQ]
```

### Installation

```bash
# Clone repository
git clone https://github.com/JordanMarcelino/learn-go-microservices.git
cd learn-go-microservices

# Set up environment (go to each services)
cp .env.example .env

# Start with Docker Compose
docker-compose up -d --build
```
