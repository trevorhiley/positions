package fakegenerator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/trevorhiley/positions/pkg/position"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//CreateFakePositions create fake positions
func CreateFakePositions(numberOfPorfolios int, numberOfInvestments int, numberOfLots int) error {
	dirName := "../../files"
	removeContents(dirName)

	jobs := make(chan int, numberOfPorfolios)
	results := make(chan int, 100)

	for w := 1; w <= 8; w++ {
		go worker(w, numberOfInvestments, numberOfLots, dirName, jobs, results)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 0; j <= numberOfPorfolios; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a <= numberOfPorfolios; a++ {
		<-results
	}

	return nil
}

func worker(id int, numberOfInvestments int, numberOfLots int, dirName string, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		//fmt.Println("worker", id, "started  job", j)
		err := createPosition(numberOfInvestments, numberOfLots, dirName)
		if err != nil {
			fmt.Printf("An error occurred: %v", err)
		}
		//fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(position position.Header, dirName string) error {

	positionsJSON, err := json.Marshal(position)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	// To start, here's how to dump a string (or just
	// bytes) into a file.
	d1 := []byte(positionsJSON)
	filename := strings.Replace(position.Portfolios[0].PortfolioID.String(), "-", "", -1)
	err = ioutil.WriteFile(fmt.Sprintf("%s/postions_%s.json", dirName, filename), d1, 0644)

	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	return nil
}

func createPosition(numberOfInvestments, numberOfLots int, dirName string) error {
	positions, err := createHeader()

	if err != nil {
		return fmt.Errorf("An error occurred")
	}
	newPortfolio, err := createPortfolio(numberOfInvestments, numberOfLots)
	if err != nil {
		return fmt.Errorf("An error occurred")
	}

	positions.Portfolios = append(positions.Portfolios, newPortfolio)

	err = writeFile(positions, dirName)

	if err != nil {
		return fmt.Errorf("An error occurred")
	}

	return nil
}

func createHeader() (position.Header, error) {
	positions := position.Header{}

	messageID, err := uuid.NewUUID()

	if err != nil {
		return positions, fmt.Errorf("An error occurred")
	}

	positions.MessageID = messageID

	positions.GeneratedDateTime = time.Now()

	return positions, nil
}

func createPortfolio(numberOfInvestments int, numberOfLots int) (position.Portfolio, error) {
	messageID, err := uuid.NewUUID()

	if err != nil {
		return position.Portfolio{}, fmt.Errorf("An error occurred")
	}

	newPortfolio := position.Portfolio{
		PortfolioID:       messageID,
		AsofDate:          time.Now(),
		AccountingBasisID: 1,
	}

	for i := 0; i <= numberOfInvestments; i++ {
		newInvestment, err := createInvestment(numberOfLots)

		if err != nil {
			return position.Portfolio{}, fmt.Errorf("An error occurred")
		}

		newPortfolio.Investments = append(newPortfolio.Investments, newInvestment)
	}

	return newPortfolio, nil
}

func createInvestment(numberOfLots int) (position.Investment, error) {
	newInvestment := position.Investment{
		InvestmentID: createRandomInt(80, 30000),
		CostValue:    createRandomFloat(100000, 100000000),
		BookValue:    createRandomFloat(100000, 100000000),
		Quantity:     createRandomFloat(10000, 1000000),
	}

	for i := 0; i <= numberOfLots; i++ {
		newLot, err := createLot()

		if err != nil {
			return position.Investment{}, fmt.Errorf("An error occurred")
		}

		newInvestment.Lots = append(newInvestment.Lots, newLot)
	}

	return newInvestment, nil
}

func createLot() (position.Lot, error) {

	messageID, err := uuid.NewUUID()

	if err != nil {
		return position.Lot{}, fmt.Errorf("An error occurred")
	}

	newLot := position.Lot{
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
