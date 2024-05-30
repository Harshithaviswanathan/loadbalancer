package main

import (
	// "flag"
	// "loadbalancer/loadbalancer"
	 "loadbalancer/servers"
)

func main() {
	// var algoStr string
	// flag.StringVar(&algoStr, "algo", "roundrobin", "Load balancing algorithm (roundrobin, weightedroundrobin, leastconnection, iphash)")
	// flag.Parse()

	// var algo loadbalancer.BalancingAlgorithm
	// switch algoStr {
	// case "roundrobin":
	// 	algo = loadbalancer.RoundRobin
	// case "weightedroundrobin":
	// 	algo = loadbalancer.WeightedRoundRobin
	// case "leastconnection":
	// 	algo = loadbalancer.LeastConnection
	// case "iphash":
	// 	algo = loadbalancer.IPHash
	// default:
	// 	panic("Unknown load balancing algorithm")
	// }

	// //Set weights for servers
	// weights := []int{1, 0, 3, 4, 2}

	// initialConnections := []int{0, 1, 0, 2, 0}

	// //Start the load balancer
	// loadbalancer.MakeLoadBalancer(5, weights, initialConnections, algo)

	 servers.RunServers(5)
}
run this as go run main.go in one terminal
package main

import (
	"flag"
	"loadbalancer/loadbalancer"
	//  "loadbalancer/servers"
)

func main() {
	var algoStr string
	flag.StringVar(&algoStr, "algo", "roundrobin", "Load balancing algorithm (roundrobin, weightedroundrobin, leastconnection, iphash)")
	flag.Parse()

	var algo loadbalancer.BalancingAlgorithm
	switch algoStr {
	case "roundrobin":
		algo = loadbalancer.RoundRobin
	case "weightedroundrobin":
		algo = loadbalancer.WeightedRoundRobin
	case "leastconnection":
		algo = loadbalancer.LeastConnection
	case "iphash":
		algo = loadbalancer.IPHash
	default:
		panic("Unknown load balancing algorithm")
	}

	//Set weights for servers
	weights := []int{1, 0, 3, 4, 2}

	initialConnections := []int{0, 1, 0, 2, 0}

	//Start the load balancer
	loadbalancer.MakeLoadBalancer(5, weights, initialConnections, algo)

	//  servers.RunServers(5)
}
run this as go run main.go in another terminal

