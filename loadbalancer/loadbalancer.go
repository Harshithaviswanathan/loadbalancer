package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LoadBalancer struct {
	RevProxy      httputil.ReverseProxy
	Endpoints     Endpoints
	BalancingAlgo BalancingAlgorithm
}

func MakeLoadBalancer(amount int, weights []int, initialConnections []int, algo BalancingAlgorithm) {
	// Instantiate Objects
	var lb LoadBalancer
	lb.BalancingAlgo = algo
	lb.Endpoints.Populate(amount, weights, initialConnections)

	// Perform initial health check
	lb.Endpoints.HealthCheck()

	// Start health check goroutine
	go lb.Endpoints.RunHealthCheck()

	// Server + Router
	router := http.NewServeMux()

	// Handler Functions
	router.HandleFunc("/loadbalancer", lb.makeRequest)

	// Listen and Serve
	log.Fatal(http.ListenAndServe(":7667", router))
}

func (lb *LoadBalancer) makeRequest(w http.ResponseWriter, r *http.Request) {
	var endpoint *url.URL
	switch lb.BalancingAlgo {
	case RoundRobin:
		endpoint = lb.Endpoints.GetNextRoundRobin()
	case WeightedRoundRobin:
		endpoint = lb.Endpoints.GetNextWeightedRoundRobin()
	case LeastConnection:
		endpoint = lb.Endpoints.GetNextLeastConnection()
	case IPHash:
		endpoint = lb.Endpoints.GetIPHash(r.RemoteAddr)
	}

	if endpoint != nil {
		lb.RevProxy = *httputil.NewSingleHostReverseProxy(endpoint)
		lb.RevProxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
	}
}
