package tsp

import (
	"math"
	"math/rand"
)

// StartStepsPar makes the ant select the next node to construct their path
func (ant *Ant) StartStepsPar(solverChan chan []Node) {
	for {
		if ant.Done {
			return
		}
		// Calculate denominator of p value
		pDenom := float64(0)
		for i := range ant.Unvisited {
			pDenom += ant.Tsp.P(&ant.CurrentNode, &ant.Unvisited[i])
		}
		// Select next node by calculating the p value
		n := ant.Rand.Float64()
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
			ant.Done = true
			ant.PathLength = ant.Tsp.Length(ant.Path)
			solverChan <- ant.Path
			break
		}
	}
}

// ParSolver runs the Ant Colony System paralelly
type ParSolver struct {
	Tsp *Tsp
	Ants []Ant
	BestPathLength float64
	BestPath []Node
}

// Solve method runs ACS paralelly
func (solver *ParSolver) Solve() {
	solver.BestPathLength = math.MaxFloat64
	solverChan := make(chan []Node)
	// Create ants
	for i := 0; i < solver.Tsp.AntCount; i++ {
		ant := Ant {ID: i, Tsp: solver.Tsp}
		ant.Rand = rand.New(rand.NewSource(int64(ant.Tsp.Rand.Intn(999999999))))
		solver.Ants = append(solver.Ants, ant)
	}
	// Run iterations
	for itr := 0; itr < solver.Tsp.ItrCount; itr++ {
		solver.Tsp.DoneAntCount = 0
		// Reset ants
		for j := range solver.Ants {
			solver.Ants[j].Reset()
		}
		for j := range solver.Ants {
			go solver.Ants[j].StartStepsPar(solverChan)
		}
		for {
			// Wait for an ant to sent its path to us
			antPath := <-solverChan
			antPathLength := solver.Tsp.Length(antPath)
			if antPathLength < solver.BestPathLength {
				solver.BestPathLength = antPathLength
				solver.BestPath = make([]Node, len(antPath))
				copy(solver.BestPath, antPath)
			}
			solver.Tsp.DoneAntCount++
			if solver.Tsp.DoneAntCount == solver.Tsp.AntCount {
				break
			}
		}
		// Let ants apply local pheromone update
		for j := range solver.Ants {
			solver.Ants[j].UpdatePheros()
		}
		// Global pheromone update (evaporation)
		for x := 0; x < len(solver.Tsp.Nodes) - 1; x++ {
			for y := x + 1; y < len(solver.Tsp.Nodes); y++ {
				xNode := &solver.Tsp.Nodes[x]
				yNode := &solver.Tsp.Nodes[y]
				solver.Tsp.SetPhero(xNode, yNode, solver.Tsp.Pheros[x][y] * (1 - solver.Tsp.Ro))
			}
		}

	}
}