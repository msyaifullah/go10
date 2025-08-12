// internal/services/interfaces.go
package services

import (
	"loan-service/internal/models"

	"github.com/google/uuid"
)

// Service interfaces
type LoanServiceInterface interface {
	GetLoanByID(id uuid.UUID) (*models.LoanSummaryResponse, error)

	// Process Loan
	ProcessCreateLoan(req *models.CreateLoanRequest) (*models.LoanSummaryResponse, error)
	ProcessApproveLoan(id uuid.UUID, req *models.CreateApprovalRequest) (*models.LoanApprovalResponse, error)
	ProcessInvestment(loanID uuid.UUID, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error)
	ProcessDisbursement(loanID uuid.UUID, req *models.CreateDisbursementRequest) (*models.DisbursementResponse, error)
}
