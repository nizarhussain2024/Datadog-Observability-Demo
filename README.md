# Datadog Observability Demo

Demonstration of Datadog observability platform integration with Docker Compose.

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Application Layer                           │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Go Application                                │  │
│  │  - HTTP endpoints                                      │  │
│  │  - Business logic                                      │  │
│  │  - Instrumented with Datadog                          │  │
│  └───────────────────────┬───────────────────────────────┘  │
└───────────────────────────┼──────────────────────────────────┘
                            │
                            │ Metrics, Traces, Logs
                            │
┌───────────────────────────▼──────────────────────────────────┐
│              Datadog Agent (Sidecar)                          │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Datadog Agent Container                        │  │
│  │  - Metrics collection                                  │  │
│  │  - Trace collection                                    │  │
│  │  - Log collection                                      │  │
│  │  - APM (Application Performance Monitoring)            │  │
│  └───────────────────────┬───────────────────────────────┘  │
└───────────────────────────┼──────────────────────────────────┘
                            │
                            │ HTTPS
                            │
┌───────────────────────────▼──────────────────────────────────┐
│              Datadog Cloud Platform                           │
│  ┌──────────────────┐  ┌──────────────────┐                 │
│  │   Metrics        │  │   APM            │                 │
│  │   Dashboard     │  │   Traces        │                 │
│  └──────────────────┘  └──────────────────┘                 │
│  ┌──────────────────┐  ┌──────────────────┐                 │
│  │   Logs           │  │   Alerts        │                 │
│  │   Management     │  │   Monitoring    │                 │
│  └──────────────────┘  └──────────────────┘                 │
└──────────────────────────────────────────────────────────────┘
```

### Observability Components

**Application:**
- Go HTTP service
- Datadog APM instrumentation
- Structured logging

**Datadog Agent:**
- Metrics collection
- Trace collection
- Log collection
- APM agent

**Datadog Platform:**
- Metrics visualization
- Distributed tracing
- Log management
- Alerting

## Design Decisions

### Observability Strategy
- **Datadog APM**: Automatic instrumentation
- **Structured Logs**: JSON format
- **Metrics**: Custom and automatic metrics
- **Traces**: Distributed tracing

### Technology Choices
- **Agent**: Datadog Agent (sidecar container)
- **Instrumentation**: Datadog APM library
- **Logging**: Structured JSON logs
- **Metrics**: StatsD protocol

### Architecture Patterns
- **Sidecar Pattern**: Datadog Agent as sidecar
- **Automatic Instrumentation**: Zero-code changes
- **Centralized Collection**: All data to Datadog

## End-to-End Flow

### Flow 1: Application Metrics Collection

```
1. Application Running
   └─> Go application starts
       └─> Datadog APM library initialized
           └─> Automatic instrumentation enabled

2. Request Processing
   └─> HTTP request received
       └─> Application processes request:
           ├─> Business logic execution
           ├─> Database queries (if any)
           └─> Response generation

3. Automatic Instrumentation
   └─> Datadog APM library:
       ├─> Captures request metadata:
       │   ├─> HTTP method
       │   ├─> URL path
       │   ├─> Status code
       │   └─> Duration
       ├─> Creates trace span
       └─> Records metrics

4. Metrics Emission
   └─> Application emits metrics:
       ├─> HTTP request count
       ├─> Request duration
       ├─> Error rate
       └─> Custom business metrics
       └─> Metrics sent via StatsD to Datadog Agent

5. Agent Collection
   └─> Datadog Agent:
       ├─> Receives metrics via StatsD (port 8125)
       ├─> Receives traces via APM (port 8126)
       ├─> Receives logs via TCP/HTTP
       └─> Aggregates and buffers data

6. Data Transmission
   └─> Datadog Agent:
       ├─> Batch metrics every 10 seconds
       ├─> Send to Datadog API (HTTPS)
       └─> Include tags:
           ├─> service: app
           ├─> environment: docker
           └─> version: 1.0.0

7. Datadog Platform
   └─> Datadog receives data:
       ├─> Store metrics in time-series DB
       ├─> Index traces
       ├─> Parse and store logs
       └─> Make available for querying

8. Visualization
   └─> User accesses Datadog dashboard:
       ├─> View metrics graphs
       ├─> Analyze traces
       └─> Search logs
```

### Flow 2: Distributed Tracing

```
1. Request Initiation
   └─> Client sends request to application
       └─> HTTP GET /api/endpoint

2. Trace Creation
   └─> Datadog APM:
       ├─> Creates root span:
       │   ├─> Trace ID: abc123
       │   ├─> Span ID: span-1
       │   └─> Operation: http.request
       └─> Starts timing

3. Application Processing
   └─> Application executes:
       ├─> Create child span: "database.query"
       ├─> Execute database query
       ├─> End child span
       └─> Create child span: "external.api.call"
           └─> Call external service
           └─> End child span

4. Trace Completion
   └─> Request completes:
       ├─> End root span
       ├─> Add tags:
       │   ├─> http.method: GET
       │   ├─> http.status_code: 200
       │   └─> http.url: /api/endpoint
       └─> Send trace to Datadog Agent

5. Agent Processing
   └─> Datadog Agent:
       ├─> Receives trace via APM port
       ├─> Validates trace data
       └─> Queues for transmission

6. Trace Transmission
   └─> Agent sends trace to Datadog:
       ├─> Batch with other traces
       ├─> Compress data
       └─> Send via HTTPS

7. Trace Storage
   └─> Datadog:
       ├─> Index trace by trace ID
       ├─> Store span data
       └─> Make searchable

8. Trace Visualization
   └─> User queries Datadog APM:
       ├─> Search by trace ID or service
       ├─> View trace timeline:
       │   ├─> Root span: http.request (50ms)
       │   ├─> Child: database.query (20ms)
       │   └─> Child: external.api.call (25ms)
       └─> Identify performance bottlenecks
```

### Flow 3: Log Collection and Analysis

```
1. Application Logging
   └─> Application emits log:
       └─> Structured JSON:
       {
         "timestamp": "2024-01-01T00:00:00Z",
         "level": "INFO",
         "message": "Request processed",
         "trace_id": "abc123",
         "service": "app",
         "method": "GET",
         "path": "/api/endpoint"
       }

2. Log Output
   └─> Application writes to stdout/stderr
       └─> Docker captures logs

3. Agent Log Collection
   └─> Datadog Agent:
       ├─> Reads from Docker logs
       ├─> Parses JSON logs
       ├─> Extracts fields
       └─> Adds metadata:
           ├─> container_id
           ├─> image_name
           └─> environment

4. Log Processing
   └─> Agent:
       ├─> Parse log format
       ├─> Extract structured fields
       ├─> Apply log processing rules
       └─> Add tags

5. Log Transmission
   └─> Agent sends logs to Datadog:
       ├─> Batch logs
       ├─> Compress
       └─> Send via HTTPS

6. Log Storage
   └─> Datadog:
       ├─> Index logs
       ├─> Store in log management system
       └─> Make searchable

7. Log Analysis
   └─> User queries Datadog Logs:
       ├─> Search: service:app status:error
       ├─> Filter by time range
       ├─> View log details
       └─> Correlate with traces (via trace_id)
```

## Data Flow

```
Application
    │
    ├─> Metrics (StatsD) ──> Datadog Agent
    ├─> Traces (APM) ──────> Datadog Agent
    └─> Logs (stdout) ──────> Datadog Agent
                              │
                              │ HTTPS
                              ▼
                         Datadog Cloud
                              │
                              ▼
                    ┌─────────┴─────────┐
                    │                   │
              Metrics DB          Logs Index
              Traces DB           APM Service
```

## Configuration

### Docker Compose
```yaml
services:
  app:
    build: .
    environment:
      - DD_AGENT_HOST=datadog
      - DD_TRACE_ENABLED=true
    depends_on:
      - datadog

  datadog:
    image: gcr.io/datadoghq/agent:latest
    environment:
      - DD_API_KEY=${DD_API_KEY}
      - DD_APM_ENABLED=true
      - DD_LOGS_ENABLED=true
    ports:
      - "8126:8126"  # APM
      - "8125:8125"  # StatsD
```

### Application Configuration
```go
// Enable Datadog APM
import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

tracer.Start(
    tracer.WithService("app"),
    tracer.WithEnv("docker"),
)
defer tracer.Stop()
```

## Build & Run

### Prerequisites
- Docker and Docker Compose
- Datadog API key

### Setup
```bash
# Set Datadog API key
export DD_API_KEY=your-api-key

# Start services
docker-compose up --build
```

### Access
- **Application**: http://localhost:8081
- **Health Check**: http://localhost:8081/health
- **Datadog Dashboard**: https://app.datadoghq.com

## Datadog Features

### Metrics
- Automatic HTTP metrics
- Custom business metrics
- Infrastructure metrics
- Real-time dashboards

### APM (Tracing)
- Distributed tracing
- Service map
- Performance insights
- Error tracking

### Logs
- Centralized log management
- Log search and filtering
- Log correlation with traces
- Log-based alerts

### Alerts
- Metric-based alerts
- Log-based alerts
- Anomaly detection
- On-call management

## Future Enhancements

- [ ] Custom metrics
- [ ] Synthetic monitoring
- [ ] Real User Monitoring (RUM)
- [ ] Network monitoring
- [ ] Security monitoring
- [ ] Cost optimization
- [ ] Custom dashboards
- [ ] Alerting rules
- [ ] Integration with CI/CD
- [ ] Performance optimization insights

## AI/NLP Capabilities

This project includes AI and NLP utilities for:
- Text processing and tokenization
- Similarity calculation
- Natural language understanding

*Last updated: 2025-12-20*
