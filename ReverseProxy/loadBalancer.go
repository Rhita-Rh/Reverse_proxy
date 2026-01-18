package main

import "net/url"

type LoadBalancer interface {
	GetNextValidPeer_RoundRobin() *Backend
	GetNextValidPeer_LeastConn() *Backend
	AddBackend(backend *Backend)
	SetBackendStatus(uri *url.URL, alive bool)
}

