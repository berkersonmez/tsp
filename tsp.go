package tsp

import (
	"math"
)

// Node represents a node such as a city
type Node struct {
	ID int
	X  float64
	Y  float64
}

// Tsp contains variables and common helper methods for solvers
type Tsp struct {
	Nodes    []Node
	Seed     int
	Alpha    float64
	Beta     float64
	Ro       float64
	AntCount int
	ItrCount int
	Dists    [][]float64
	Pheros   [][]float64
	Q        float64
}

// Distance returns distance between nodes n1 and n2
func (tsp *Tsp) Distance(n1, n2 *Node) float64 {
	return math.Sqrt(math.Pow(n1.X-n2.X, 2) + math.Pow(n1.Y-n2.Y, 2))
}

// AddNode adds a node to the tsp struct
func (tsp *Tsp) AddNode(n1 *Node) {
	tsp.Nodes = append(tsp.Nodes, *n1)
}

// SetPhero sets the pheromone value for the edge
func (tsp *Tsp) SetPhero(n1, n2 *Node, phero float64) {
	tsp.Pheros[n1.ID][n2.ID] = phero
	tsp.Pheros[n2.ID][n1.ID] = phero
}

// UpdatePhero increments the pheromone by value of delta
func (tsp *Tsp) UpdatePhero(n1, n2 *Node, delta float64) {
	tsp.SetPhero(n1, n2, tsp.Pheros[n1.ID][n2.ID])
}

// Initialize sets initial pheromones and distances
func (tsp *Tsp) Initialize() {
	initPhero := float64(1) / float64(len(tsp.Nodes))
	minDist := math.MaxFloat64
	tsp.Dists = make([][]float64, len(tsp.Nodes))
	tsp.Pheros = make([][]float64, len(tsp.Nodes))
	for i := range tsp.Nodes {
		tsp.Dists[i] = make([]float64, len(tsp.Nodes))
		tsp.Pheros[i] = make([]float64, len(tsp.Nodes))

		for j := range tsp.Dists[i] {
			dist := tsp.Distance(&tsp.Nodes[i], &tsp.Nodes[j])
			tsp.Dists[i][j] = dist
			tsp.Dists[j][i] = dist
			tsp.Pheros[i][j] = initPhero
			tsp.Pheros[j][i] = initPhero
			if dist < minDist {
				minDist = dist
			}
		}
	}
	tsp.Q = minDist
}

// Tau gets the tau value for the edge
func (tsp *Tsp) Tau(n1, n2 *Node) float64 {
	return tsp.Pheros[n1.ID][n2.ID]
}

// Nu gets the nu value for the edge
func (tsp *Tsp) Nu(n1, n2 *Node) float64 {
	return 1 / tsp.Dists[n1.ID][n2.ID]
}

// P calculates the p value for the edge
func (tsp *Tsp) P(n1, n2 *Node) float64 {
	return math.Pow(tsp.Tau(n1, n2), tsp.Alpha) * math.Pow(tsp.Nu(n1, n2), tsp.Beta)
}