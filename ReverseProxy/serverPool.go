package main

import (
	"net/url"
)

type ServerPool struct{
	Backends[ ]*Backend `json:"backends"`
	Current uint64 `json:"current"`
}

func (serverPool *ServerPool) GetNextValidPeer_RoundRobin() *Backend{
	backends := serverPool.Backends
	if len(backends)==0{
		return nil
	}

	for i :=0; i<len(backends); i++{
		nextValidPeer_index := (int(serverPool.Current) + i) % len(backends)
		backend := backends[nextValidPeer_index]
		if backend.Alive{
			serverPool.Current = uint64(nextValidPeer_index +1)
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
		if !backend.Alive{
			continue
		}
		if min_conn_backend == nil || backend.CurrentConns < min_conn_backend.CurrentConns{
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
			backend.Alive = alive
			return 
		}
	}
}
