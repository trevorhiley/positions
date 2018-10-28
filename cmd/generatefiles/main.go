package main

import (
	"fmt"

	"github.com/trevorhiley/positions/pkg/fakegenerator"
)

const numberOfPortfolios int = 10
const numberOfInvestments int = 10
const numberOfLots int = 100

//Main runs the package
func main() {
	_, err := fakegenerator.CreateFakePositions(numberOfPortfolios, numberOfInvestments, numberOfLots)

	if err != nil {
		fmt.Print(err)
	}

	//fmt.Print(positions)

}
