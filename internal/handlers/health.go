package handlers

import (
	"net/http"
	"sync"
)

var (
	healthyMutex sync.RWMutex
	healthy      = false
)

func UpdateHealthy(isHealthy bool) {
	healthyMutex.Lock()
	healthy = isHealthy
	healthyMutex.Unlock()
}

func HealthCheckEndpoint(resp http.ResponseWriter, _ *http.Request) {
	healthyMutex.RLock()
	isHealthy := healthy
	healthyMutex.RUnlock()
	if isHealthy {
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("OK"))
	} else {
		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write([]byte("Unhealthy"))
	}
}
