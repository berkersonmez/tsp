package tsp

import (
	"math"
	"math/rand"
	"errors"
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
	Seed     int64
	Alpha    float64
	Beta     float64
	Ro       float64
	AntCount int
	ItrCount int
	Dists    [][]float64
	Pheros   [][]float64
	Q        float64
	Rand     *rand.Rand
	DoneAntCount int
}

// Distance returns distance between nodes n1 and n2
func (tsp *Tsp) Distance(n1, n2 *Node) float64 {
	return math.Sqrt(math.Pow(n1.X-n2.X, 2) + math.Pow(n1.Y-n2.Y, 2))
}

// Length returns the length of the path
func (tsp *Tsp) Length(nodes []Node) float64 {
	current := &nodes[0]
	totalLength := float64(0)
	for i := range nodes {
		prev := current
		current = &nodes[i]
		totalLength += tsp.Distance(prev, current)
	}
	totalLength += tsp.Distance(&nodes[0], current)
	return totalLength
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
	tsp.SetPhero(n1, n2, tsp.Pheros[n1.ID][n2.ID] + delta)
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
	}
	for i := 0; i < len(tsp.Nodes) - 1; i++ {
		for j := i + 1; j < len(tsp.Nodes); j++ {
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
	tsp.Rand = rand.New(rand.NewSource(tsp.Seed))
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

// Solver is the interface for the main Ant Colony System algorithm
type Solver interface {
	Solve()
}

// === Some general functions below ===

// IndexOf finds the index of node in the slice
func IndexOf(n1 *Node, slc []Node) (int, error) {
	for i := range slc {
		if slc[i].ID == n1.ID {
			return i, nil
		}
	}
	return 0, errors.New("Node not in slice")
}

// Remove returns a slice with the designated element removed
func Remove(n1 *Node, slc []Node) []Node {
	indexToRemove, _ := IndexOf(n1, slc)
	return append(slc[:indexToRemove], slc[indexToRemove+1:]...)
}

// Ant is the builder of the solution in the Ant Colony System
type Ant struct {
	ID          int
	Path        []Node
	Unvisited   []Node
	CurrentNode Node
	Tsp         *Tsp
	PathLength  float64
	Done        bool
}

// Reset reinitializes the ant with an empty path
func (ant *Ant) Reset() {
	ant.Unvisited = make([]Node, len(ant.Tsp.Nodes))
	copy(ant.Unvisited, ant.Tsp.Nodes)
	ant.Path = make([]Node, 0)
	ant.CurrentNode = ant.Unvisited[ant.Tsp.Rand.Intn(len(ant.Unvisited))]
	ant.Done = false
	ant.PathLength = 0
	ant.Path = append(ant.Path, ant.CurrentNode)
	ant.Unvisited = Remove(&ant.CurrentNode, ant.Unvisited)
}

// Step makes the ant select the next node to construct their path
func (ant *Ant) Step() {
	if ant.Done {
		return
	}
	// Calculate denominator of p value
	pDenom := float64(0)
	for i := range ant.Unvisited {
		pDenom += ant.Tsp.P(&ant.CurrentNode, &ant.Unvisited[i])
	}
	// Select next node by calculating the p value
	n := ant.Tsp.Rand.Float64()
	movingSum := float64(0)
	for i := range ant.Unvisited {
		p := ant.Tsp.P(&ant.CurrentNode, &ant.Unvisited[i]) / pDenom
		movingSum += p
		if movingSum >= n {
			ant.Path = append(ant.Path, ant.Unvisited[i])
			ant.CurrentNode = ant.Unvisited[i]
			ant.Unvisited = Remove(&ant.Unvisited[i], ant.Unvisited)
			break
		}
	}
	// Check if ant is done
	if (len(ant.Unvisited) == 0) {
		ant.Tsp.DoneAntCount++
		ant.Done = true
		ant.PathLength = ant.Tsp.Length(ant.Path)
	}

}

// UpdatePheros updates pheromones according to ant's path
func (ant *Ant) UpdatePheros() {
	// Local pheromone update
	deltaTau := ant.Tsp.Q / ant.PathLength
	x := &ant.Path[0]
	for i := range ant.Path[1:] {
		y := &ant.Path[1:][i]
		ant.Tsp.UpdatePhero(x, y, deltaTau)
		x = y
	}
}