package models

import (
	"time"

	"github.com/google/uuid"
)

// CreateLoanRequest represents the request to create a new loan
type CreateLoanRequest struct {
	BorrowerID      uuid.UUID `json:"borrower_id" validate:"required"`
	PrincipalAmount float64   `json:"principal_amount" validate:"required,gt=0"`
	InterestRate    float64   `json:"interest_rate" validate:"required,gte=0,lte=1"`
	ROI             float64   `json:"roi" validate:"required,gte=0,lte=1"`
}

// UpdateLoanStateRequest represents the request to update loan state
type UpdateLoanStateRequest struct {
	NewState     LoanState `json:"new_state" validate:"required"`
	ChangedBy    uuid.UUID `json:"changed_by" validate:"required"`
	ChangeReason string    `json:"change_reason,omitempty"`
}

// CreateApprovalRequest represents the request to approve a loan
type CreateApprovalRequest struct {
	LoanID              uuid.UUID `json:"loan_id" validate:"required"`
	ValidatorID         uuid.UUID `json:"validator_id" validate:"required"`
	ApprovalDate        time.Time `json:"approval_date" validate:"required"`
	VisitProofImageURL  string    `json:"visit_proof_image_url" validate:"required"`
	VisitProofImageType FileType  `json:"visit_proof_image_type" validate:"required"`
	Notes               string    `json:"notes,omitempty"`
}

// CreateInvestmentRequest represents the request to create an investment
type CreateInvestmentRequest struct {
	LoanID         uuid.UUID `json:"loan_id" validate:"required"`
	InvestorID     uuid.UUID `json:"investor_id" validate:"required"`
	Amount         float64   `json:"amount" validate:"required,gt=0"`
	InvestmentDate time.Time `json:"investment_date" validate:"required"`
}

// CreateDisbursementRequest represents the request to disburse a loan
type CreateDisbursementRequest struct {
	LoanID                  uuid.UUID `json:"loan_id" validate:"required"`
	FieldOfficerID          uuid.UUID `json:"field_officer_id" validate:"required"`
	DisbursementDate        time.Time `json:"disbursement_date" validate:"required"`
	SignedAgreementURL      string    `json:"signed_agreement_url" validate:"required"`
	SignedAgreementFileType FileType  `json:"signed_agreement_file_type" validate:"required"`
	DisbursedAmount         float64   `json:"disbursed_amount" validate:"required,gt=0"`
	Notes                   string    `json:"notes,omitempty"`
}

// CreateBorrowerRequest represents the request to create a new borrower
type CreateBorrowerRequest struct {
	IDNumber    string `json:"id_number" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address"`
}

// UpdateBorrowerRequest represents the request to update borrower information
type UpdateBorrowerRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Address     string `json:"address,omitempty"`
}

// CreateEmployeeRequest represents the request to create a new employee
type CreateEmployeeRequest struct {
	EmployeeID  string `json:"employee_id" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Role        string `json:"role" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateEmployeeRequest represents the request to update employee information
type UpdateEmployeeRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	Role        string `json:"role,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// CreateInvestorRequest represents the request to create a new investor
type CreateInvestorRequest struct {
	InvestorCode string `json:"investor_code" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	PhoneNumber  string `json:"phone_number"`
}

// UpdateInvestorRequest represents the request to update investor information
type UpdateInvestorRequest struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
