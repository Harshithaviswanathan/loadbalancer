package loadbalancer

const (
	baseURL = "http://localhost:808"
)

type BalancingAlgorithm int

const (
	RoundRobin BalancingAlgorithm = iota
	WeightedRoundRobin
	LeastConnection
	IPHash
)
