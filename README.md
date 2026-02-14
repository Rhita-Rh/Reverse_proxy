# Reverse Proxy

A lightweight, high-performance reverse proxy written in Go with load balancing capabilities, health checks, and an admin API.

## Overview

This project implements a reverse proxy server that distributes incoming HTTP requests across multiple backend servers using load balancing strategies. It includes built-in health checking to ensure requests are only routed to healthy backends.

## Features

- **Load Balancing**: round-robin and least connections
- **Health Checks**: Automatic health monitoring of backend servers every 30 sec
- **Admin API**: Administrative endpoints for monitoring and configuration
- **Configurable**: JSON-based configuration for easy setup
- **Multiple Backend Support**: Route requests to multiple backend servers
- **Lightweight**: Built with Go standard library for minimal dependencies

## Project Structure

```
Reverse_proxy/
├── config.json              # Configuration file
├── go.mod                   # Go module definition
├── README.md                # This file
├── Client/
│   └── client.go           # Test client for sending requests
├── ReverseProxy/
│   ├── reverseProxy.go     # Main reverse proxy implementation
│   ├── loadBalancer.go     # Load balancing interface
│   ├── serverPool.go       # Backend server pool management
│   ├── adminApi.go         # Admin API endpoints
│   └── Backend.go          # Backend server representation
├── Server1/
│   └── server1.go          # Backend server 1
├── Server2/
│   └── server2.go          # Backend server 2
└── Server3/
    └── server3.go          # Backend server 3
```

## Configuration

The proxy is configured via `config.json`:

```json
{
  "proxy": {
    "port": 8080,
    "strategy": "round-robin",
    "health_check_frequency": "30s"
  },
  "backends": [
    {
      "url": "http://localhost:8081"
    },
    {
      "url": "http://localhost:8082"
    },
    {
      "url": "http://localhost:8083"
    }
  ]
}
```

### Configuration Options

- **port**: The port where the reverse proxy listens (default: 8080)
- **strategy**: Load balancing strategy - `"round-robin"` or `"least-conn"` (default: round-robin)
- **health_check_frequency**: How often to check backend health (default: 30s)
- **backends**: Array of backend server URLs to route traffic to

## Getting Started

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd reverse_proxy
```



### Running the Project

1. **Start Backend Servers** (in separate terminals):

```bash
# Terminal 1
go run Server1/server1.go

# Terminal 2
go run Server2/server2.go

# Terminal 3
go run Server3/server3.go
```
2. **Start the Reverse Proxy**:
```bash
go run ReverseProxy/*.go

or

go run ReverseProxy/reverseProxy.go ReverseProxy/backend.go ReverseProxy/serverPool.go ReverseProxy/adminApi.go
```

The proxy will start listening on `http://localhost:8080`

3. **Test with the Client**:(in separate terminals)

```bash
go run Client1/client1.go
```

if you want to test the lean-conn method run these as well to test:

```bash
go run Client2/client2.go

go run Client3/client3.go
```

## Usage

### Sending Requests

Once the reverse proxy is running, send HTTP requests to `http://localhost:8080`. The proxy will:

1. Receive the request
2. Select a backend using the configured load balancing strategy (you can change it in config.json "round-robin" or "least-conn")
3. Forward the request to the selected backend
4. Return the response to the client

### Load Balancing Strategies

- **Round-Robin**: Distributes requests evenly across all healthy backends
- **Least Connections**: Routes requests to the backend with the fewest active connections

### Health Checks

The proxy automatically monitors backend server health at the frequency specified in the configuration(30sec).

## API Endpoints

### Admin API

The admin API provides endpoints for monitoring and managing the proxy (if implemented in `adminApi.go`): (run these commands in cmd and not powershell)

● GET /status: Return a JSON list of all backends and their current health/load:
```bash
curl http://localhost:9000/status
```


● POST /backends: Dynamically add a new backend URL to the pool:
```bash
curl -X POST http://localhost:9000/backends -H "Content-Type: application/json" -d "{\"url\": \"http://localhost:8084\"}"
```


● DELETE /backends: Remove a backend from the pool, example:(server running on port :8083)
```bash
curl -X DELETE http://localhost:9000/backends -H "Content-Type: application/json" -d "{\"url\": \"http://localhost:8083\"}"
```

## Architecture

### Components

- **ReverseProxy**: Main proxy handler that accepts incoming requests
- **ServerPool**: Manages the collection of backend servers and their health status
- **LoadBalancer**: Implements load balancing strategies
- **Backend**: Represents individual backend servers
- **AdminApi**: Provides administrative functionality

## Development

## Performance Considerations

- Configurable health check frequency to balance between responsiveness and overhead
- Efficient round-robin algorithm for minimal latency
- Go's concurrency model for handling multiple simultaneous requests

## Troubleshooting

### Port Already in Use
If you get a "port already in use" error, either:
- Change the port in `config.json`
- Don't use port `:8080` because it's used by the proxy
- Don't use port `:9090` because it's used by the adminApi

### Backends Not Responding
- Ensure all backend servers are running on their configured ports
- Verify URLs in `config.json`

### Health Check Issues
- Increase `health_check_frequency` if experiencing high load

## Future Enhancements

- [ ] TLS/HTTPS support
- [ ] Cookie-based session affinity
- [ ] Request rate limiting
- [ ] Weighted round-robin
- [ ] Graceful shutdown

