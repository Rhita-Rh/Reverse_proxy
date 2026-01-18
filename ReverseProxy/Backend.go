package main

import (
	"net/url"
	"sync"
)

type Backend struct{
	URL *url.URL `json:"-"`
	URLString string `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
	mux sync.RWMutex
}

func (backend *Backend)IncrementConn(){
	backend.CurrentConns ++
}

func (backend *Backend)DecrementConn(){
	backend.CurrentConns --
}