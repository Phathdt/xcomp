# Asynq Scheduler and Processor Setup

This document describes the asynq-based async job processing system implemented in this project.

## Overview

The async system consists of:
- **CheckPendingOrderScheduler**: Enqueues jobs every 5 seconds
- **CheckPendingOrderProcessor**: Processes jobs to log pending order and customer information
- **Asynq Monitor**: Web interface for job monitoring

## Components

### 1. Job Types (`jobs/types.go`)
- `CheckPendingOrderJob`: Job payload structure
- `TypeCheckPendingOrder`: Job type constant

### 2. Scheduler (`schedulers/check_pending_order_scheduler.go`)
- **CheckPendingOrderScheduler**: Schedules jobs every 5 seconds
- Enqueues `CheckPendingOrderJob` to the default queue
- Handles graceful shutdown

### 3. Processor (`processors/check_pending_order_processor.go`)
- **CheckPendingOrderProcessor**: Processes pending order jobs
- Fetches all pending orders from the order service
- Logs detailed information about each pending order and associated customer
- Includes order items details in debug logs

### 4. Async Module (`infrastructure/async/async.module.go`)
- **AsyncService**: Coordinates scheduler and processor
- Sets up asynq server with queue configuration
- Configures asynq monitor for job monitoring
- Handles service lifecycle (start/stop)

## Configuration

### Redis Configuration
```yaml
redis:
  host: 'localhost'
  port: 6379
  password: ''
  db: 0
```

### Async Monitor Configuration
```yaml
async:
  monitor:
    port: 8080
    enabled: true  # false in production
```

## Queue Configuration

The system uses three queues with different priorities:
- **critical**: 6 workers (highest priority)
- **default**: 3 workers (medium priority)
- **low**: 1 worker (lowest priority)

## Monitoring

Asynq Monitor is available at:
- Development: `http://localhost:8080/monitoring`
- Port configurable via `async.monitor.port` in config

## Usage

### Starting the System
The async system starts automatically when the main application starts:
```bash
go run . serve
```

### Viewing Logs
The processor logs detailed information about pending orders:
```
INFO  Pending order found order_id=xxx customer_username=john_doe customer_email=john@example.com
DEBUG Order item details product_name="Product A" quantity=2 unit_price=29.99
```

### Monitoring Jobs
Access the monitoring interface at the configured port to:
- View job queues and their status
- Monitor job processing rates
- Inspect failed jobs
- Retry failed jobs

## Dependencies

- `github.com/hibiken/asynq`: Main async job processing library
- `github.com/hibiken/asynqmon`: Web-based monitoring interface
- Redis: Required for job queue storage

## Security Notes

- Monitoring interface should be disabled in production (`async.monitor.enabled: false`)
- Redis should be properly secured with authentication
- Consider using Redis Sentinel or Cluster for high availability
