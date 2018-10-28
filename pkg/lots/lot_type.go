package lots

import (
	"github.com/google/uuid"
)

//Lot struct type for storing position lot information
type Lot struct {
	LotID       uuid.UUID
	BookValue   float64
	MarketValue float64
	CostValue   float64
	Quantity    float64
}
