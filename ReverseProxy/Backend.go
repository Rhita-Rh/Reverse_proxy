package main

import (
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct{
	URL *url.URL `json:"-"`
	URLString string `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
	mux sync.RWMutex
}

type BackendConfig struct{
	URLString string `json:"url"`
}

// using atomic for concurrency 
func (backend *Backend)IncrementConn(){
	atomic.AddInt64(&backend.CurrentConns, 1)
}

func (backend *Backend)DecrementConn(){
	atomic.AddInt64(&backend.CurrentConns, -1)
}

func (backend *Backend) IsAlive() bool{
	backend.mux.Lock()
	defer backend.mux.Unlock()
	return backend.Alive
}

func (backend *Backend) GetCurentConn()int64{
	return atomic.LoadInt64(&backend.CurrentConns)
}