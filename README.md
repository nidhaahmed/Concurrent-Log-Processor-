# ⚔️ Concurrent Log Processing System

A high-performance log processing backend built in Go with concurrency, modular architecture, and microservice design principles.

---

## 🌐 Live Demo

🔗 https://concurrent-log-processor.onrender.com/

Upload a log file and compare sequential vs concurrent processing performance in real-time.

## 🚀 Features

- 📂 Log ingestion via REST API
- ⚡ Concurrent processing using worker pools (goroutines + channels)
- 🧠 Pluggable parser architecture (Simple + Regex)
- 📊 Aggregation of log levels (INFO, ERROR, WARNING)
- 🔍 Filtering and querying logs
- 🧾 In-memory storage layer
- 📦 Dockerized deployment
- 📈 Metrics & request logging
- 🌐 UI for log upload + performance comparison

---

## 🧠 System Architecture

```
Client → API → Processor → Parser → Aggregator → Storage
```

---

## ⚡ Concurrency Model

- Worker pool using goroutines
- Jobs distributed via channels
- Thread-safe aggregation using mutex
- Benchmark comparison (sequential vs concurrent)

---

## 🖥️ UI Demo

Upload a log file and compare:

- Sequential processing time
- Concurrent processing time
- Aggregated results

---

## 📡 API Endpoints

| Endpoint | Method | Description |
|--------|--------|-------------|
| /health | GET | Health check |
| /upload | POST | Process logs |
| /logs | GET | Get all logs |
| /filter | GET | Filter logs by level |
| /compare | POST | Compare sequential vs concurrent |
| /metrics | GET | Request metrics |

---

## 🐳 Docker Usage

### Build Image
```bash
docker build -t log-processor .
```

### Run Container
```bash
docker run -p 8080:8080 log-processor
```

---

## 🧪 Example Request

```json
{
  "logs": [
    "2026-04-27 ERROR Database failed",
    "2026-04-27 INFO Server started"
  ]
}
```

---

## 📊 Sample Output

```json
{
  "concurrent": {
    "time_ms": 5,
    "result": {
      "ERROR": 1,
      "INFO": 1
    }
  },
  "sequential": {
    "time_ms": 15,
    "result": {
      "ERROR": 1,
      "INFO": 1
    }
  }
}
```

---

## 🛠️ Tech Stack

- Go (Golang)
- net/http
- Goroutines & Channels
- Docker
- JSON APIs

---

## 💡 Key Learnings

- Designing concurrent systems using worker pools
- Avoiding bottlenecks in channel pipelines
- Thread-safe data structures using mutex
- Building modular and scalable backend architecture
- Dockerizing backend services for portability

---

## 🚀 Future Improvements

- Persistent storage (Redis / DB)
- Kafka-based log ingestion
- Distributed worker nodes
- Advanced analytics dashboard

---

## 👨‍💻 Author

Nidha Ahmed Mohammad
