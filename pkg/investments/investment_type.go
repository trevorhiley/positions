package investments

import "github.com/trevorhiley/positions/pkg/lots"

//Investment struct type for storing position investment information
type Investment struct {
	InvestmentID  int
	PositionPrice float64
	Quantity      float64
	BookValue     float64
	MarketValue   float64
	CostValue     float64
	Lots          []lots.Lot
}
