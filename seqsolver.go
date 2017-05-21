package tsp

import (
	"math"
)

// SeqSolver runs the Ant Colony System sequentally
type SeqSolver struct {
	Tsp *Tsp
	Ants []Ant
	BestPathLength float64
	BestPath []Node
}

// Solve method runs ACS sequentally
func (solver *SeqSolver) Solve() {
	solver.BestPathLength = math.MaxFloat64
	// Create ants
	for i := 0; i < solver.Tsp.AntCount; i++ {
		ant := Ant {ID: i, Tsp: solver.Tsp}
		solver.Ants = append(solver.Ants, ant)
	}

	// Run iterations
	for itr := 0; itr < solver.Tsp.ItrCount; itr++ {
		solver.Tsp.DoneAntCount = 0
		// Reset ants
		for j := range solver.Ants {
			solver.Ants[j].Reset()
		}
		// Let ants construct their paths
		for solver.Tsp.DoneAntCount < len(solver.Ants) {
			for j := range solver.Ants {
				solver.Ants[j].Step()
			}
		}
		// Let ants apply local pheromone update
		for j := range solver.Ants {
			solver.Ants[j].UpdatePheros()
		}
		// Find global (including previous iters) best path among the ant paths
		for j := range solver.Ants {
			if solver.Ants[j].PathLength < solver.BestPathLength {
				solver.BestPathLength = solver.Ants[j].PathLength
				solver.BestPath = make([]Node, len(solver.Ants[j].Path))
				copy(solver.BestPath, solver.Ants[j].Path)
			}
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