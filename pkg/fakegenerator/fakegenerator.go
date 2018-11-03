package fakegenerator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/google/uuid"
	"github.com/trevorhiley/positions/pkg/position"
)

var svc *s3.S3

func init() {
	rand.Seed(time.Now().UnixNano())
}

//CreateFakePositions create fake positions
func CreateFakePositions(numberOfPorfolios int, numberOfInvestments int, numberOfLots int) error {
	svc = s3.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	dirName := "../../files"
	removeContents(dirName)

	jobs := make(chan int, numberOfPorfolios)
	results := make(chan int, 100)

	for w := 1; w <= 8; w++ {
		go worker(w, numberOfInvestments, numberOfLots, dirName, jobs, results)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 0; j < numberOfPorfolios; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a < numberOfPorfolios; a++ {
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

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(positionsJSON); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	// To start, here's how to dump a string (or just
	// bytes) into a file.
	dateString := position.Portfolios[0].AsofDate.Format("2006-01-02")
	uuidNoDash := position.Portfolios[0].PortfolioID.String()
	filename := fmt.Sprintf("/%s/%s/%s_positions.json", dateString, uuidNoDash, uuidNoDash)
	//err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", dirName, filename), d1, 0644)

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(b.Bytes())),
		Bucket: aws.String("positions-sample-files-go"),
		Key:    aws.String(filename),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(err)
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err)
			return err
		}
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
		AsofDate:          createAsofDate(time.Now()),
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

func createAsofDate(timeIn time.Time) time.Time {
	year, month, day := timeIn.Date()
	parsedAsofDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return parsedAsofDate
}

func createRandomInt(min int, max int) int {

	return (rand.Intn(max - min)) + min
}

func createRandomFloat(min int, max int) float64 {
	return (rand.Float64()) * float64(createRandomInt(min, max))
}
