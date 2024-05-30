package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Endpoint struct {
	URL       *url.URL
	Healthy   bool
	Mutex     sync.Mutex
	NumActive int // Number of active connections
}

type Endpoints struct {
	List            []*Endpoint
	RoundRobinIdx   int
	ConnectionMap   map[string]int // For LeastConnection
	WeightedCounts  map[string]int // For WeightedRoundRobin
	Mutex           sync.Mutex
	requestCounter  int   // Counter for tracking requests
	originalWeights []int // Original weights of servers
}

func (e *Endpoints) Populate(amount int, weights []int, initialConnections []int) {
	fmt.Println("Populating endpoints...")
	for i := 0; i < amount; i++ {
		fmt.Printf("Endpoint %d: Initial Connections: %d\n", i, initialConnections[i])
		e.List = append(e.List, &Endpoint{
			URL:       createEndpoint(baseURL, i),
			Healthy:   true,
			NumActive: initialConnections[i], // Set initial number of active connections
		})
	}

	// Initialize the connection map for LeastConnection
	e.ConnectionMap = make(map[string]int)

	// Initialize the weighted counts for WeightedRoundRobin
	e.WeightedCounts = make(map[string]int)
	for i, endpoint := range e.List {
		e.ConnectionMap[endpoint.URL.String()] = initialConnections[i]
		e.WeightedCounts[endpoint.URL.String()] = weights[i]
	}

	// Save the original weights
	e.originalWeights = make([]int, len(weights))
	copy(e.originalWeights, weights)
}

// HealthCheck performs health checks on endpoints
func (e *Endpoints) HealthCheck() {
	for _, endpoint := range e.List {
		resp, err := http.Get(endpoint.URL.String())
		if err != nil || resp.StatusCode != http.StatusOK {
			endpoint.Mutex.Lock()
			endpoint.Healthy = false
			endpoint.Mutex.Unlock()
			log.Printf("Endpoint %s is unhealthy: %v", endpoint.URL.String(), err)
		} else {
			endpoint.Mutex.Lock()
			endpoint.Healthy = true
			endpoint.Mutex.Unlock()
			log.Printf("Endpoint %s is healthy!", endpoint.URL.String())
		}
	}
}

// RunHealthCheck starts a routine for periodic health checks
func (e *Endpoints) RunHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			e.HealthCheck()
		}
	}
}

// GetNextRoundRobin returns the next endpoint using round-robin algorithm
func (e *Endpoints) GetNextRoundRobin() *url.URL {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()
	for i := 0; i < len(e.List); i++ {
		endpoint := e.List[e.RoundRobinIdx]
		e.RoundRobinIdx = (e.RoundRobinIdx + 1) % len(e.List)
		if endpoint.Healthy {
			return endpoint.URL
		}
	}
	// Fallback if no healthy endpoint found
	return nil
}

func (e *Endpoints) GetNextLeastConnection() *url.URL {
	minConn := -1
	var minEndpoint *url.URL

	e.Mutex.Lock()
	defer e.Mutex.Unlock()

	for _, endpoint := range e.List {
		if endpoint.Healthy {
			connCount := e.ConnectionMap[endpoint.URL.String()] // Get current connection count
			fmt.Printf("Endpoint: %s, Connections: %d\n", endpoint.URL.String(), connCount)
			if minConn == -1 || connCount < minConn {
				minConn = connCount
				minEndpoint = endpoint.URL
			}
		}
	}

	if minEndpoint != nil {
		// Increment connection count for the selected endpoint
		e.ConnectionMap[minEndpoint.String()]++
	}

	fmt.Printf("Selected Endpoint: %s, Connections: %d\n", minEndpoint.String(), minConn)

	return minEndpoint
}

func (e *Endpoints) GetNextWeightedRoundRobin() *url.URL {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()

	// Select the endpoint based on weighted round-robin
	for i := 0; i < len(e.List); i++ {
		idx := (e.RoundRobinIdx + i) % len(e.List)
		endpoint := e.List[idx]
		if endpoint.Healthy && e.WeightedCounts[endpoint.URL.String()] > 0 {
			e.RoundRobinIdx = (idx + 1) % len(e.List)
			e.WeightedCounts[endpoint.URL.String()]--
			e.requestCounter++
			if e.requestCounter%10 == 0 {
				// Reset the weights after every 10th request
				e.resetWeights()
			}
			return endpoint.URL
		}
	}

	return nil
}

func (e *Endpoints) resetWeights() {
	for i, endpoint := range e.List {
		e.WeightedCounts[endpoint.URL.String()] = e.originalWeights[i]
	}
}

func (e *Endpoints) GetIPHash(remoteAddr string) *url.URL {
	// Use a consistent hash function for the client's remote address
	hash := hashFunction(remoteAddr)
	initialIdx := int(hash % uint64(len(e.List))) // Convert to int here
	idx := initialIdx
	// Start searching from the initially selected index
	for {
		// Check if the selected server is healthy
		if e.List[idx].Healthy {
			return e.List[idx].URL
		}
		// Move to the next server in the list
		idx = (idx + 1) % len(e.List)
		// If we have searched all servers and none are healthy, return nil
		if idx == initialIdx {
			return nil
		}
	}
}

// Example of a simple consistent hash function
func hashFunction(input string) uint64 {
	var hash uint64 = 5381
	for _, c := range input {
		hash = ((hash << 5) + hash) + uint64(c)
	}
	return hash
}

func createEndpoint(endpoint string, idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	url, _ := url.Parse(link)
	return url
}
