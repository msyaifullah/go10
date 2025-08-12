// internal/repositories/interfaces.go
package repositories

import (
	"database/sql"
	"loan-service/internal/models"
	"time"

	"github.com/google/uuid"
)

// Repository interface
type LoanRepositoryInterface interface {

	// Loan creating and basic operations
	CreateLoan(tx *sql.Tx, loan *models.Loan) (*models.Loan, error)
	GetLoanByID(tx *sql.Tx, loanID uuid.UUID) (*models.Loan, error)
	UpdateLoanState(tx *sql.Tx, loanID uuid.UUID, newState models.LoanState) (*models.Loan, error)
	RecordLoanStateHistory(tx *sql.Tx, prevState models.LoanState, loan *models.Loan, employeeID uuid.UUID, changeReason string) (*models.LoanStateHistory, error)

	// loan approval (Proposed → Approved)
	CreateApproval(tx *sql.Tx, approval *models.Approval) (*models.Approval, error)

	// loan investment (Approved → Invested)
	CreateInvestment(tx *sql.Tx, investment *models.Investment) (*models.Investment, error)
	UpdateLoanTotalInvested(tx *sql.Tx, loanID uuid.UUID, newTotal float64) (*models.Loan, error)
	GetInvestmentsNeedingAgreementEmail() ([]*models.Investment, error)
	UpdateInvestmentAgreementSent(investmentID uuid.UUID, agreementSent bool, agreementSentAt *time.Time) error
	UpdateLoanAgreementLetterURL(tx *sql.Tx, loanID uuid.UUID, agreementURL string) error

	// loan disbursement (Invested → Disbursed)
	CreateDisbursement(tx *sql.Tx, disbursement *models.Disbursement) (*models.Disbursement, error)

	// User Management
	GetInvestorByID(investorID uuid.UUID) (*models.Investor, error)
	GetBorrowerByID(borrowerID uuid.UUID) (*models.Borrower, error)

	// Communication & Notifications
	CreateEmailNotification(notification *models.EmailNotification) error
}
