# Distributed Job Scheduling Platform

## Overview
Distributed job scheduling platform supports both one-time and recurring cron jobs, with jobs organized either as independent tasks or as Directed Acyclic Graphs (DAGs). The platform ensures scalability, fault tolerance, and observability, leveraging a microservices-based architecture.

### Core Features:
- **Priority-based Queue Management**
- **Dead Letter Queue Handling**
- **Job Retry Mechanisms with Exponential Backoff**
- **Real-time Job Status Monitoring**
- **Horizontal Scaling**
- **Prometheus/Grafana Integration for Monitoring**

---

## Architecture

### Workflow:
1. **User**:
   - Initiates job submission via the User Interface or API Gateway.

2. **API Gateway**:
   - Routes user requests to appropriate microservices.

3. **Job Submission Service**:
   - Validates and submits jobs to Kafka.
   - Jobs are prioritized and added to a queue.

4. **Scheduler Service**:
   - Fetches jobs from Kafka and schedules them based on execution time.
   - Handles retries for failed jobs.

5. **Job Executor Service**:
   - Executes scheduled jobs and updates the Job Status Service.

6. **Job Status Service**:
   - Tracks the current status of all jobs (Pending, In-Progress, Completed, Failed).

7. **DAG Orchestration Service**:
   - Manages jobs with dependencies and ensures DAG execution order.


## Upcoming features:

8. **Monitoring Service**:
   - Exposes metrics via Prometheus for system health and performance.
   - Dashboards are available in Grafana for real-time insights.

9. **Notification Service**:
   - Sends alerts and notifications (email, SMS, etc.) based on job events.

10. **User Interface**:
    - Provides users with a dashboard to track job statuses and manage tasks.

---

## Microservices

### 1. **Job Submission Service**
- Validates and enqueues jobs into Kafka.
- Manages job priority.

### 2. **Scheduler Service**
- Schedules jobs for execution.
- Implements retry logic with exponential backoff.

### 3. **Job Executor Service**
- Executes jobs and logs results.
- Communicates with the Job Status Service.

### 4. **Job Status Service**
- Maintains real-time job statuses.
- Exposes APIs for querying job statuses.

### 5. **DAG Orchestration Service**
- Ensures correct execution of dependent tasks within DAGs.

### 6. **Monitoring Service**
- Uses Prometheus for metrics collection.
- Integrates with Grafana for visualization.
---

## Installation and Setup

### Prerequisites
- Go 1.20+
- Kafka
- Docker (for containerized deployment)
- Prometheus and Grafana

### Steps:
1. Clone the repository:
   ```bash
   git clone https://github.com/avalokitasharma/Distributed-Task-Scheduling.git
   cd distributed-job-scheduling
   ```
2. Build and run individual services:
   ```bash
   cd services/job-submission
   go build -o job-submission
   ./job-submission
   ```
3. Start Kafka and dependencies using Docker Compose:
   ```bash
   docker-compose up -d
   ```
4. Set up Prometheus and Grafana for monitoring:
   - Edit the `prometheus.yml` configuration file.
   - Import Grafana dashboards.

---

## Usage

### Submitting a Job
- Use the API Gateway:
  ```bash
  curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
        "name": "example-job",
        "type": "cron",
        "schedule": "0 * * * *",
        "priority": 1
     }'
  ```

### Viewing Job Status
- Query the Job Status Service:
  ```bash
  curl http://localhost:8080/api/job-status/{job-id}
  ```

### Monitoring
- Access Grafana at [http://localhost:3000](http://localhost:3000).
- Prometheus metrics are available at [http://localhost:9090](http://localhost:9090).

---

## Contribution Guidelines
1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/new-feature
   ```
3. Commit changes and push to your fork.
4. Open a pull request.

---
