package fakegenerator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/trevorhiley/positions/pkg/header"
	"github.com/trevorhiley/positions/pkg/investments"
	"github.com/trevorhiley/positions/pkg/lots"
	"github.com/trevorhiley/positions/pkg/portfolios"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//CreateFakePositions create fake positions
func CreateFakePositions(numberOfPorfolios int, numberOfInvestments int, numberOfLots int) (header.Header, error) {
	positions, err := createHeader()

	if err != nil {
		return positions, fmt.Errorf("An error occurred")
	}

	for i := 0; i <= numberOfPorfolios; i++ {
		newPortfolio, err := createPortfolio(numberOfInvestments, numberOfLots)
		if err != nil {
			return positions, fmt.Errorf("An error occurred")
		}

		positions.Portfolios = append(positions.Portfolios, newPortfolio)
	}

	return positions, nil
}

func createHeader() (header.Header, error) {
	positions := header.Header{}

	messageID, err := uuid.NewUUID()

	if err != nil {
		return positions, fmt.Errorf("An error occurred")
	}

	positions.MessageID = messageID

	positions.GeneratedDateTime = time.Now()

	return positions, nil
}

func createPortfolio(numberOfInvestments int, numberOfLots int) (portfolios.Portfolio, error) {
	messageID, err := uuid.NewUUID()

	if err != nil {
		return portfolios.Portfolio{}, fmt.Errorf("An error occurred")
	}

	newPortfolio := portfolios.Portfolio{
		PortfolioID:       messageID,
		AsofDate:          time.Now(),
		AccountingBasisID: 1,
	}

	for i := 0; i <= numberOfInvestments; i++ {
		newInvestment, err := createInvestment(numberOfLots)

		if err != nil {
			return portfolios.Portfolio{}, fmt.Errorf("An error occurred")
		}

		newPortfolio.Investments = append(newPortfolio.Investments, newInvestment)
	}

	return newPortfolio, nil
}

func createInvestment(numberOfLots int) (investments.Investment, error) {
	newInvestment := investments.Investment{
		InvestmentID: createRandomInt(80, 30000),
		CostValue:    createRandomFloat(100000, 100000000),
		BookValue:    createRandomFloat(100000, 100000000),
		Quantity:     createRandomFloat(10000, 1000000),
	}

	for i := 0; i <= numberOfLots; i++ {
		newLot, err := createLot()

		if err != nil {
			return investments.Investment{}, fmt.Errorf("An error occurred")
		}

		newInvestment.Lots = append(newInvestment.Lots, newLot)
	}

	return newInvestment, nil
}

func createLot() (lots.Lot, error) {

	messageID, err := uuid.NewUUID()

	if err != nil {
		return lots.Lot{}, fmt.Errorf("An error occurred")
	}

	newLot := lots.Lot{
		LotID:     messageID,
		CostValue: createRandomFloat(100000, 100000000),
		BookValue: createRandomFloat(100000, 100000000),
		Quantity:  createRandomFloat(10000, 1000000),
	}

	return newLot, nil
}

func createRandomInt(min int, max int) int {

	return (rand.Intn(max - min)) + min
}

func createRandomFloat(min int, max int) float64 {
	return (rand.Float64()) * float64(createRandomInt(min, max))
}
