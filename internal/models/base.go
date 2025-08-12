package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// LoanState represents the current state of a loan
type LoanState string

const (
	LoanStateProposed  LoanState = "proposed"
	LoanStateApproved  LoanState = "approved"
	LoanStateInvested  LoanState = "invested"
	LoanStateDisbursed LoanState = "disbursed"
)

// Valid states for transition validation
var validStates = map[LoanState]bool{
	LoanStateProposed:  true,
	LoanStateApproved:  true,
	LoanStateInvested:  true,
	LoanStateDisbursed: true,
}

// Valid state transitions
var validTransitions = map[LoanState][]LoanState{
	LoanStateProposed:  {LoanStateApproved},
	LoanStateApproved:  {LoanStateInvested},
	LoanStateInvested:  {LoanStateDisbursed},
	LoanStateDisbursed: {}, // Final state, no transitions allowed
}

// IsValid checks if the loan state is valid
func (ls LoanState) IsValid() bool {
	return validStates[ls]
}

// CanTransitionTo checks if transition from current state to target state is valid
func (ls LoanState) CanTransitionTo(target LoanState) bool {
	allowedStates := validTransitions[ls]
	for _, state := range allowedStates {
		if state == target {
			return true
		}
	}
	return false
}

// String returns the string representation of the loan state
func (ls LoanState) String() string {
	return string(ls)
}

// FileType represents supported file types for uploads
type FileType string

const (
	FileTypePDF  FileType = "pdf"
	FileTypeJPEG FileType = "jpeg"
	FileTypePNG  FileType = "png"
)
