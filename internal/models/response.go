package models

import (
	"time"

	"github.com/google/uuid"
)

// LoanSummaryResponse represents a summary view of a loan
type LoanSummaryResponse struct {
	ID                  uuid.UUID  `json:"id"`
	BorrowerName        string     `json:"borrower_name"`
	PrincipalAmount     float64    `json:"principal_amount"`
	InterestRate        float64    `json:"interest_rate"`
	ROI                 float64    `json:"roi"`
	State               LoanState  `json:"state"`
	TotalInvested       float64    `json:"total_invested"`
	RemainingInvestment float64    `json:"remaining_investment"`
	InvestorCount       int        `json:"investor_count"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	ApprovalDate        *time.Time `json:"approval_date,omitempty"`
	DisbursementDate    *time.Time `json:"disbursement_date,omitempty"`
}

type LoanApprovalResponse struct {
	ID                  uuid.UUID `json:"id"`
	LoanID              uuid.UUID `json:"loan_id"`
	ValidatorID         uuid.UUID `json:"validator_id"`
	ApprovalDate        time.Time `json:"approval_date"`
	VisitProofImageURL  string    `json:"visit_proof_image_url"`
	VisitProofImageType FileType  `json:"visit_proof_image_type"`
	Notes               string    `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// BorrowerSummaryResponse represents a summary view of a borrower
type BorrowerSummaryResponse struct {
	ID          uuid.UUID `json:"id"`
	IDNumber    string    `json:"id_number"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	LoanCount   int       `json:"loan_count"`
	TotalLoaned float64   `json:"total_loaned"`
	CreatedAt   time.Time `json:"created_at"`
}

// InvestorSummaryResponse represents a summary view of an investor
type InvestorSummaryResponse struct {
	ID              uuid.UUID `json:"id"`
	InvestorCode    string    `json:"investor_code"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	IsActive        bool      `json:"is_active"`
	InvestmentCount int       `json:"investment_count"`
	TotalInvested   float64   `json:"total_invested"`
	ExpectedReturns float64   `json:"expected_returns"`
	CreatedAt       time.Time `json:"created_at"`
}

// EmployeeSummaryResponse represents a summary view of an employee
type EmployeeSummaryResponse struct {
	ID          uuid.UUID `json:"id"`
	EmployeeID  string    `json:"employee_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// InvestmentResponse represents the response for investment creation
type InvestmentResponse struct {
	ID             uuid.UUID `json:"id"`
	LoanID         uuid.UUID `json:"loan_id"`
	InvestorID     uuid.UUID `json:"investor_id"`
	Amount         float64   `json:"amount"`
	ExpectedReturn float64   `json:"expected_return"`
	InvestmentDate time.Time `json:"investment_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// // InvestmentSummaryResponse represents a summary view of an investment
// type InvestmentSummaryResponse struct {
// 	ID             uuid.UUID `json:"id"`
// 	LoanID         uuid.UUID `json:"loan_id"`
// 	InvestorID     uuid.UUID `json:"investor_id"`
// 	InvestorName   string    `json:"investor_name"`
// 	BorrowerName   string    `json:"borrower_name"`
// 	Amount         float64   `json:"amount"`
// 	ExpectedReturn float64   `json:"expected_return"`
// 	InvestmentDate time.Time `json:"investment_date"`
// 	AgreementSent  bool      `json:"agreement_sent"`
// 	CreatedAt      time.Time `json:"created_at"`
// }

// DisbursementResponse represents the response for disbursement creation
type DisbursementResponse struct {
	ID                      uuid.UUID `json:"id"`
	LoanID                  uuid.UUID `json:"loan_id"`
	FieldOfficerID          uuid.UUID `json:"field_officer_id"`
	DisbursementDate        time.Time `json:"disbursement_date"`
	SignedAgreementURL      string    `json:"signed_agreement_url"`
	SignedAgreementFileType FileType  `json:"signed_agreement_file_type"`
	DisbursedAmount         float64   `json:"disbursed_amount"`
	Notes                   string    `json:"notes"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
