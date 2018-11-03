package main

import (
	"fmt"

	"github.com/trevorhiley/positions/pkg/fakegenerator"
)

const numberOfPortfolios int = 1000
const numberOfInvestments int = 10
const numberOfLots int = 100

//Main runs the package
func main() {

	//previousmax := runtime.GOMAXPROCS(1)
	//fmt.Printf("previous max prox was %v", previousmax)

	err := fakegenerator.CreateFakePositions(numberOfPortfolios, numberOfInvestments, numberOfLots)

	if err != nil {
		fmt.Print(err)
	}

}
