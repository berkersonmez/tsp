package tsp

import (
	"math"
	"fmt"
)

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
			fmt.Printf("DEBUG: Path length=%v\n",solver.Ants[j].PathLength)
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