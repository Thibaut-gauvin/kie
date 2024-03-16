package handlers

import (
	"net/http"
	"sync"
)

var (
	readyMutex sync.RWMutex
	ready      = false
)

func UpdateReady(isReady bool) {
	readyMutex.Lock()
	ready = isReady
	readyMutex.Unlock()
}

func ReadyCheckEndpoint(resp http.ResponseWriter, _ *http.Request) {
	readyMutex.RLock()
	isReady := ready
	readyMutex.RUnlock()
	if isReady {
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("OK"))
	} else {
		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write([]byte("Unhealthy"))
	}
}
