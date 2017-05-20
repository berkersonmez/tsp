package main

import (
	"os"
	"strconv"
	"bufio"
	"strings"
	"fmt"
	"errors"
	"github.com/berkersonmez/tsp"
)

func main() {
	var vtsp tsp.Tsp
	var err error
	if len(os.Args) < 7 {
		check(errors.New("Please specify 7 command line arguments (filename, seed, alpha, beta, ro, antcount, itercount)"))
	}
	filename := os.Args[1]
	vtsp.Seed, err = strconv.ParseInt(os.Args[2], 10, 64)
	check(err)
	vtsp.Alpha, err = strconv.ParseFloat(os.Args[3], 64)
	check(err)
	vtsp.Beta, err = strconv.ParseFloat(os.Args[4], 64)
	check(err)
	vtsp.Ro, err = strconv.ParseFloat(os.Args[5], 64)
	check(err)
	vtsp.AntCount, err = strconv.Atoi(os.Args[6])
	check(err)
	vtsp.ItrCount, err = strconv.Atoi(os.Args[7])
	check(err)

	fmt.Println("* Reading file " + filename + "...")
	file, err := os.Open(filename)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	i := 0
	for scanner.Scan() {
		coords := strings.Fields(scanner.Text())
		node := tsp.Node{}
		node.ID = i
		node.X, _ = strconv.ParseFloat(coords[0], 64)
		node.Y, _ = strconv.ParseFloat(coords[1], 64)
		vtsp.AddNode(&node)
		i++
	}
	file.Close()
	fmt.Println("* Successfully read file.")

	fmt.Println("* Solving TSP sequentally...")
	vtsp.Initialize()
	seqSolver := tsp.SeqSolver{Tsp: &vtsp}
	seqSolver.Solve()
	fmt.Println("* Sequental solution done.")
	fmt.Printf("Result: %v\n", seqSolver.BestPathLength)

}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}