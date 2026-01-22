package main

import (
	"net/url"
	"sync"
)

type ServerPool struct{
	Backends[ ]*Backend `json:"backends"`
	Current uint64 `json:"current"`
	mux sync.Mutex
}

func (serverPool *ServerPool) GetNextValidPeer_RoundRobin() *Backend{
	serverPool.mux.Lock()
	defer serverPool.mux.Unlock()

	n := len(serverPool.Backends)
	backends := serverPool.Backends
	if n == 0{
		return nil
	}

	startIdx := serverPool.Current

	for i :=0; i<n; i++{
		nextValidPeer_index := (int(startIdx) + i) % n
		backend := backends[nextValidPeer_index]
		if backend.IsAlive(){
			serverPool.Current = uint64(nextValidPeer_index+1) % uint64(n)
			return backend
		}
	}
	return nil // If all backends are down
}

func (serverPool *ServerPool) GetNextValidPeer_LeastConn() *Backend{
	backends := serverPool.Backends
	if len(backends) == 0{
		return nil
	}
	var min_conn_backend *Backend
	for _, backend := range backends{
		if !backend.IsAlive(){
			continue
		}
		if min_conn_backend == nil || backend.GetCurentConn() < min_conn_backend.GetCurentConn(){
			min_conn_backend = backend
		}
	}
	
	return min_conn_backend
	
}

func (serverPool *ServerPool) AddBackend(backend *Backend){
	serverPool.Backends = append(serverPool.Backends, backend)
}

func (serverPool *ServerPool) SetBackendStatus(uri *url.URL, alive bool){
	backends := serverPool.Backends
	for _, backend := range backends{
		if backend.URLString == uri.String(){
			backend.mux.Lock()
			backend.Alive = alive
			backend.mux.Unlock()
			return 
		}
	}
}
