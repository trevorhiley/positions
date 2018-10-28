package portfolios

import (
	"time"

	"github.com/google/uuid"
	"github.com/trevorhiley/positions/pkg/investments"
)

//Portfolio struct type for storing portfolio lot information
type Portfolio struct {
	PortfolioID       uuid.UUID
	AsofDate          time.Time
	AccountingBasisID int
	Investments       []investments.Investment
}
