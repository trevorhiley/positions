package position

import (
	"time"

	"github.com/google/uuid"
)

//Header struct type for storing position header information
type Header struct {
	GeneratedDateTime time.Time
	MessageID         uuid.UUID
	Portfolios        []Portfolio
}

//Portfolio struct type for storing portfolio lot information
type Portfolio struct {
	PortfolioID       uuid.UUID
	AsofDate          time.Time
	AccountingBasisID int
	Investments       []Investment
}

//Investment struct type for storing position investment information
type Investment struct {
	InvestmentID  int
	PositionPrice float64
	Quantity      float64
	BookValue     float64
	MarketValue   float64
	CostValue     float64
	Lots          []Lot
}

//Lot struct type for storing position lot information
type Lot struct {
	LotID       uuid.UUID
	BookValue   float64
	MarketValue float64
	CostValue   float64
	Quantity    float64
}
