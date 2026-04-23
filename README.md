# Distributed Job Scheduler (Multi-Tenant, Airflow-like)

This project is a distributed job scheduling system inspired by Airflow, designed to handle multi-tenant workloads at scale. It supports job creation, scheduling, execution, and tenant-aware isolation using a policy-driven architecture.

The system separates control plane (tenant and policy management) from data plane (job scheduling and execution), enabling scalable and flexible infrastructure as new tenants and workloads are onboarded.

---

## Architecture Overview

The system consists of the following services:

- **API Service**  
  Accepts job creation requests. Injects tenant context using middleware and persists job definitions.

- **Planner Service**  
  Expands scheduled jobs into executable job runs and pushes them into tenant-aware queues.

- **Worker Service**  
  Pulls jobs from queues and executes them. Uses tenant policy to determine resource allocation.

- **Tenant Service**  
  Manages tenants, plans, and isolation policies. Acts as the control plane for multi-tenancy.

---

## Key Concepts

### Tenant-Aware Design

All requests and background processing are scoped by tenant. A tenant policy defines:

- Isolation mode (shared vs isolated)
- Rate limits
- Concurrency limits
- Resource namespaces (e.g., Redis keys)

This policy is fetched once and propagated through request context, avoiding scattered conditional logic.

---

### Policy-Driven Infrastructure

Instead of hardcoding behavior, services rely on:

- **Policy Provider** – fetches tenant configuration
- **Middleware** – injects tenant policy into request context
- **Factory Layer** – creates tenant-specific resources (queues, compute)

This allows adding new tenants, scaling infrastructure, or introducing isolation without changing service logic.

---

### Resource Abstraction

All tenant-specific behavior is abstracted behind interfaces:

- Queue (Redis-backed)
- Compute execution (worker pool)
- Storage (Postgres)

The factory layer resolves which implementation to use based on tenant policy.

---

## Data Model

Core tables:

- `tenants` – tenant metadata and limits
- `jobs` – job templates
- `job_runs` – execution instances

Each record is scoped by `tenant_id`.

---

## Running the System

### Prerequisites

- Docker
- Docker Compose

### Start services

```bash
docker-compose up --build
```
