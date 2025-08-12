package models

import (
	"time"

	"github.com/google/uuid"
)

// Investment represents individual investment in a loan
type Investment struct {
	BaseModel
	LoanID          uuid.UUID  `json:"loan_id" validate:"required"`
	InvestorID      uuid.UUID  `json:"investor_id" validate:"required"`
	Amount          float64    `json:"amount" validate:"required,gt=0"`
	InvestmentDate  time.Time  `json:"investment_date" validate:"required"`
	ExpectedReturn  float64    `json:"expected_return"`
	AgreementSent   bool       `json:"agreement_sent"`
	AgreementSentAt *time.Time `json:"agreement_sent_at,omitempty"`

	// Relationships
	Loan     *Loan     `json:"loan,omitempty"`
	Investor *Investor `json:"investor,omitempty"`
}

// CalculateExpectedReturn calculates the expected return for this investment
func (i *Investment) CalculateExpectedReturn(loanROI float64) float64 {
	return i.Amount * loanROI
}
