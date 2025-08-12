package models

import (
	"fmt"
)

// Borrower represents a loan borrower
type Borrower struct {
	BaseModel
	IDNumber    string `json:"id_number" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address"`

	// Relationships
	Loans []Loan `json:"loans,omitempty"`
}

// FullName returns the borrower's full name
func (b *Borrower) FullName() string {
	return fmt.Sprintf("%s %s", b.FirstName, b.LastName)
}

// Employee represents staff members (field validators, officers)
type Employee struct {
	BaseModel
	EmployeeID  string `json:"employee_id" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Role        string `json:"role" validate:"required"` // field_validator, field_officer, etc.
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`

	// Relationships
	Approvals     []Approval     `json:"approvals,omitempty"`
	Disbursements []Disbursement `json:"disbursements,omitempty"`
}

// FullName returns the employee's full name
func (e *Employee) FullName() string {
	return fmt.Sprintf("%s %s", e.FirstName, e.LastName)
}

// Investor represents loan investors/lenders
type Investor struct {
	BaseModel
	InvestorCode string `json:"investor_code" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	PhoneNumber  string `json:"phone_number"`
	IsActive     bool   `json:"is_active"`

	// Relationships
	Investments []Investment `json:"investments,omitempty"`
}
