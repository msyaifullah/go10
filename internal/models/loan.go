package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Loan represents the main loan entity
type Loan struct {
	BaseModel
	BorrowerID         uuid.UUID `json:"borrower_id" validate:"required"`
	PrincipalAmount    float64   `json:"principal_amount" validate:"required,gt=0"`
	InterestRate       float64   `json:"interest_rate" validate:"required,gte=0,lte=1"` // Percentage as decimal (0.10 for 10%)
	ROI                float64   `json:"roi" validate:"required,gte=0,lte=1"`           // Return on Investment for investors
	State              LoanState `json:"state" validate:"required"`
	AgreementLetterURL string    `json:"agreement_letter_url"`
	TotalInvested      float64   `json:"total_invested"`

	// Relationships - these will be populated by joins or separate queries
	Borrower     *Borrower          `json:"borrower,omitempty"`
	Approval     *Approval          `json:"approval,omitempty"`
	Investments  []Investment       `json:"investments,omitempty"`
	Disbursement *Disbursement      `json:"disbursement,omitempty"`
	StateHistory []LoanStateHistory `json:"state_history,omitempty"`
}

// ValidateStateTransition validates if the loan can transition to the target state
// Returns error with detailed reason if validation fails
func (l *Loan) ValidateStateTransition(targetState LoanState) error {
	// Check if target state is valid
	if !targetState.IsValid() {
		return fmt.Errorf("invalid target state: %s", targetState)
	}

	// Check if current state can transition to target state
	if !l.State.CanTransitionTo(targetState) {
		return fmt.Errorf("cannot transition from %s to %s", l.State, targetState)
	}

	// Business logic validation based on target state
	switch targetState {
	case LoanStateApproved:
		return l.validateApprovalTransition()
	case LoanStateInvested:
		return l.validateInvestmentTransition()
	case LoanStateDisbursed:
		return l.validateDisbursementTransition()
	}

	return nil
}

// validateApprovalTransition validates business rules for approval
func (l *Loan) validateApprovalTransition() error {
	if l.State != LoanStateProposed {
		return fmt.Errorf("loan must be in proposed state to be approved, current state: %s", l.State)
	}

	// Add any additional business rules for approval
	// For example: check if borrower exists, validate loan amount, etc.

	return nil
}

// validateInvestmentTransition validates business rules for investment
func (l *Loan) validateInvestmentTransition() error {
	if l.State != LoanStateApproved {
		return fmt.Errorf("loan must be in approved state to receive investments, current state: %s", l.State)
	}

	// Add any additional business rules for investment
	// For example: check if approval exists, validate investment amount, etc.

	return nil
}

// validateDisbursementTransition validates business rules for disbursement
func (l *Loan) validateDisbursementTransition() error {
	if l.State != LoanStateInvested {
		return fmt.Errorf("loan must be in invested state to be disbursed, current state: %s", l.State)
	}

	// Check if loan is fully invested
	if !l.IsFullyInvested() {
		return fmt.Errorf("loan must be fully invested before disbursement. Required: %.2f, Invested: %.2f",
			l.PrincipalAmount, l.TotalInvested)
	}

	// check if agreement letter is signed
	if l.AgreementLetterURL == "" {
		return fmt.Errorf("loan must have an agreement letter to be disbursed")
	}

	return nil
}

// ValidateInvestmentAmount validates if the investment amount is valid
func (l *Loan) ValidateInvestmentAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("investment amount must be greater than 0")
	}

	if l.State != LoanStateApproved {
		return fmt.Errorf("loan must be in approved state to receive investments, current state: %s", l.State)
	}

	remaining := l.RemainingInvestmentAmount()
	if amount > remaining {
		return fmt.Errorf("investment amount %.2f exceeds remaining amount %.2f", amount, remaining)
	}

	return nil
}

// RemainingInvestmentAmount calculates how much more can be invested
func (l *Loan) RemainingInvestmentAmount() float64 {
	return l.PrincipalAmount - l.TotalInvested
}

// IsFullyInvested checks if loan is fully invested
func (l *Loan) IsFullyInvested() bool {
	return l.TotalInvested >= l.PrincipalAmount
}

// LoanStateHistory tracks state changes for audit purposes
type LoanStateHistory struct {
	BaseModel
	LoanID        uuid.UUID `json:"loan_id" validate:"required"`
	PreviousState LoanState `json:"previous_state"`
	NewState      LoanState `json:"new_state" validate:"required"`
	ChangedBy     uuid.UUID `json:"changed_by" validate:"required"` // Employee who made the change
	ChangeReason  string    `json:"change_reason"`
	ChangeDate    time.Time `json:"change_date"`

	// Relationships
	Loan *Loan `json:"loan,omitempty"`
}

// Approval represents loan approval details
type Approval struct {
	BaseModel
	LoanID              uuid.UUID `json:"loan_id" validate:"required"`
	ValidatorID         uuid.UUID `json:"validator_id" validate:"required"`
	ApprovalDate        time.Time `json:"approval_date" validate:"required"`
	VisitProofImageURL  string    `json:"visit_proof_image_url" validate:"required"`
	VisitProofImageType FileType  `json:"visit_proof_image_type" validate:"required"`
	Notes               string    `json:"notes"`

	// Relationships
	Loan      *Loan     `json:"loan,omitempty"`
	Validator *Employee `json:"validator,omitempty"`
}

// Disbursement represents loan disbursement details
type Disbursement struct {
	BaseModel
	LoanID                  uuid.UUID `json:"loan_id" validate:"required"`
	FieldOfficerID          uuid.UUID `json:"field_officer_id" validate:"required"`
	DisbursementDate        time.Time `json:"disbursement_date" validate:"required"`
	SignedAgreementURL      string    `json:"signed_agreement_url" validate:"required"`
	SignedAgreementFileType FileType  `json:"signed_agreement_file_type" validate:"required"`
	DisbursedAmount         float64   `json:"disbursed_amount" validate:"required,gt=0"`
	Notes                   string    `json:"notes"`

	// Relationships
	Loan         *Loan     `json:"loan,omitempty"`
	FieldOfficer *Employee `json:"field_officer,omitempty"`
}
