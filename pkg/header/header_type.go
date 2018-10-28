package header

import (
	"time"

	"github.com/trevorhiley/positions/pkg/portfolios"

	"github.com/google/uuid"
)

//Header struct type for storing position header information
type Header struct {
	GeneratedDateTime time.Time
	MessageID         uuid.UUID
	Portfolios        []portfolios.Portfolio
}
