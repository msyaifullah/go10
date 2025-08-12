package repositories

import (
	"database/sql"
	"fmt"
	"loan-service/internal/models"
	"loan-service/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type LoanRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewLoanRepository(db *sql.DB, logger *logger.Logger) LoanRepositoryInterface {
	return &LoanRepository{
		db:     db,
		logger: logger,
	}
}

func (r *LoanRepository) CreateLoan(tx *sql.Tx, loan *models.Loan) (*models.Loan, error) {
	// Generate UUID for new loan
	loan.ID = uuid.New()

	query := `INSERT INTO loans (id, borrower_id, principal_amount, interest_rate, roi, state, agreement_letter_url, total_invested, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING created_at, updated_at`

	var err error
	if tx != nil {
		err = tx.QueryRow(query,
			loan.ID,
			loan.BorrowerID,
			loan.PrincipalAmount,
			loan.InterestRate,
			loan.ROI,
			loan.State,
			loan.AgreementLetterURL,
			loan.TotalInvested,
		).Scan(&loan.CreatedAt, &loan.UpdatedAt)
	} else {
		err = r.db.QueryRow(query,
			loan.ID,
			loan.BorrowerID,
			loan.PrincipalAmount,
			loan.InterestRate,
			loan.ROI,
			loan.State,
			loan.AgreementLetterURL,
			loan.TotalInvested,
		).Scan(&loan.CreatedAt, &loan.UpdatedAt)
	}

	return loan, err
}

func (r *LoanRepository) GetLoanByID(tx *sql.Tx, loanID uuid.UUID) (*models.Loan, error) {
	query := `
		SELECT 
			l.id, l.borrower_id, l.principal_amount, l.interest_rate, l.roi, l.state, 
			l.agreement_letter_url, l.total_invested, l.created_at, l.updated_at,
			b.id, b.id_number, b.first_name, b.last_name, b.email, b.phone_number, b.address,
			b.created_at, b.updated_at
		FROM loans l
		INNER JOIN borrowers b ON l.borrower_id = b.id
		WHERE l.id = $1 AND l.deleted_at IS NULL
	`

	var loan models.Loan
	var borrower models.Borrower

	var err error
	if tx != nil {
		err = tx.QueryRow(query, loanID).Scan(
			&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.InterestRate, &loan.ROI, &loan.State,
			&loan.AgreementLetterURL, &loan.TotalInvested, &loan.CreatedAt, &loan.UpdatedAt,
			&borrower.ID, &borrower.IDNumber, &borrower.FirstName, &borrower.LastName, &borrower.Email, &borrower.PhoneNumber, &borrower.Address,
			&borrower.CreatedAt, &borrower.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, loanID).Scan(
			&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.InterestRate, &loan.ROI, &loan.State,
			&loan.AgreementLetterURL, &loan.TotalInvested, &loan.CreatedAt, &loan.UpdatedAt,
			&borrower.ID, &borrower.IDNumber, &borrower.FirstName, &borrower.LastName, &borrower.Email, &borrower.PhoneNumber, &borrower.Address,
			&borrower.CreatedAt, &borrower.UpdatedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	loan.Borrower = &borrower
	return &loan, nil
}

func (r *LoanRepository) CreateApproval(tx *sql.Tx, approval *models.Approval) (*models.Approval, error) {
	query := `INSERT INTO approvals (id, loan_id, validator_id, approval_date, visit_proof_image_url, visit_proof_image_type, notes, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING created_at, updated_at`

	var err error
	if tx != nil {
		err = tx.QueryRow(query,
			approval.ID,
			approval.LoanID,
			approval.ValidatorID,
			approval.ApprovalDate,
			approval.VisitProofImageURL,
			approval.VisitProofImageType,
			approval.Notes,
		).Scan(&approval.CreatedAt, &approval.UpdatedAt)
	} else {
		err = r.db.QueryRow(query,
			approval.ID,
			approval.LoanID,
			approval.ValidatorID,
			approval.ApprovalDate,
			approval.VisitProofImageURL,
			approval.VisitProofImageType,
			approval.Notes,
		).Scan(&approval.CreatedAt, &approval.UpdatedAt)
	}

	return approval, err
}

func (r *LoanRepository) UpdateLoanState(tx *sql.Tx, loanID uuid.UUID, newState models.LoanState) (*models.Loan, error) {
	query := `UPDATE loans SET state = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND deleted_at IS NULL
			  RETURNING id, borrower_id, principal_amount, interest_rate, roi, state, agreement_letter_url, total_invested, created_at, updated_at`

	var loan models.Loan
	var err error
	if tx != nil {
		err = tx.QueryRow(query, newState, loanID).Scan(
			&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.InterestRate, &loan.ROI, &loan.State,
			&loan.AgreementLetterURL, &loan.TotalInvested, &loan.CreatedAt, &loan.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, newState, loanID).Scan(
			&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.InterestRate, &loan.ROI, &loan.State,
			&loan.AgreementLetterURL, &loan.TotalInvested, &loan.CreatedAt, &loan.UpdatedAt,
		)
	}

	return &loan, err
}

func (r *LoanRepository) RecordLoanStateHistory(tx *sql.Tx, prevState models.LoanState, loan *models.Loan, employeeID uuid.UUID, changeReason string) (*models.LoanStateHistory, error) {
	history := &models.LoanStateHistory{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		LoanID:        loan.ID,
		PreviousState: prevState,
		NewState:      loan.State,
		ChangedBy:     employeeID,
		ChangeReason:  changeReason,
		ChangeDate:    time.Now(),
	}

	query := `INSERT INTO loan_state_histories (id, loan_id, previous_state, new_state, changed_by, change_reason, change_date, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING created_at, updated_at`

	var dbErr error
	if tx != nil {
		dbErr = tx.QueryRow(query,
			history.ID,
			history.LoanID,
			history.PreviousState,
			history.NewState,
			history.ChangedBy,
			history.ChangeReason,
			history.ChangeDate,
		).Scan(&history.CreatedAt, &history.UpdatedAt)
	} else {
		dbErr = r.db.QueryRow(query,
			history.ID,
			history.LoanID,
			history.PreviousState,
			history.NewState,
			history.ChangedBy,
			history.ChangeReason,
			history.ChangeDate,
		).Scan(&history.CreatedAt, &history.UpdatedAt)
	}

	return history, dbErr
}

func (r *LoanRepository) CreateInvestment(tx *sql.Tx, investment *models.Investment) (*models.Investment, error) {
	query := `INSERT INTO investments (id, loan_id, investor_id, amount, expected_return, investment_date, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING created_at, updated_at`

	var err error
	if tx != nil {
		err = tx.QueryRow(query,
			investment.ID,
			investment.LoanID,
			investment.InvestorID,
			investment.Amount,
			investment.ExpectedReturn,
			investment.InvestmentDate,
		).Scan(&investment.CreatedAt, &investment.UpdatedAt)
	} else {
		err = r.db.QueryRow(query,
			investment.ID,
			investment.LoanID,
			investment.InvestorID,
			investment.Amount,
			investment.ExpectedReturn,
			investment.InvestmentDate,
		).Scan(&investment.CreatedAt, &investment.UpdatedAt)
	}

	return investment, err
}

func (r *LoanRepository) UpdateLoanTotalInvested(tx *sql.Tx, loanID uuid.UUID, newTotal float64) (*models.Loan, error) {
	// Use atomic update to prevent race conditions
	// This combines validation and update in one database operation
	query := `UPDATE loans 
		SET total_invested = total_invested + $1, 
		    updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 
		  AND deleted_at IS NULL
		  AND state = 'approved'
		  AND (total_invested + $1) <= principal_amount
		RETURNING id, borrower_id, principal_amount, interest_rate, roi, state, 
		          agreement_letter_url, total_invested, created_at, updated_at`

	var loan models.Loan
	var err error
	if tx != nil {
		err = tx.QueryRow(query, newTotal, loanID).Scan(
			&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.InterestRate, &loan.ROI, &loan.State,
			&loan.AgreementLetterURL, &loan.TotalInvested, &loan.CreatedAt, &loan.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, newTotal, loanID).Scan(
			&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.InterestRate, &loan.ROI, &loan.State,
			&loan.AgreementLetterURL, &loan.TotalInvested, &loan.CreatedAt, &loan.UpdatedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan not found, already fully invested, or not in approved state")
		}
		return nil, err
	}

	return &loan, nil
}

func (r *LoanRepository) CreateDisbursement(tx *sql.Tx, disbursement *models.Disbursement) (*models.Disbursement, error) {
	query := `INSERT INTO disbursements (id, loan_id, field_officer_id, disbursement_date, signed_agreement_url, signed_agreement_file_type, disbursed_amount, notes, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING created_at, updated_at`

	var err error
	if tx != nil {
		err = tx.QueryRow(query,
			disbursement.ID,
			disbursement.LoanID,
			disbursement.FieldOfficerID,
			disbursement.DisbursementDate,
			disbursement.SignedAgreementURL,
			disbursement.SignedAgreementFileType,
			disbursement.DisbursedAmount,
			disbursement.Notes,
		).Scan(&disbursement.CreatedAt, &disbursement.UpdatedAt)
	} else {
		err = r.db.QueryRow(query,
			disbursement.ID,
			disbursement.LoanID,
			disbursement.FieldOfficerID,
			disbursement.DisbursementDate,
			disbursement.SignedAgreementURL,
			disbursement.SignedAgreementFileType,
			disbursement.DisbursedAmount,
			disbursement.Notes,
		).Scan(&disbursement.CreatedAt, &disbursement.UpdatedAt)
	}

	return disbursement, err
}

// GetInvestmentsNeedingAgreementEmail gets investments that need agreement emails sent
func (r *LoanRepository) GetInvestmentsNeedingAgreementEmail() ([]*models.Investment, error) {
	query := `
		SELECT i.id, i.loan_id, i.investor_id, i.amount, i.investment_date, i.expected_return, 
		       i.agreement_sent, i.agreement_sent_at, i.created_at, i.updated_at
		FROM investments i
		INNER JOIN loans l ON i.loan_id = l.id AND l.deleted_at IS NULL
		WHERE i.agreement_sent = false 
		AND l.state = 'invested'
		ORDER BY i.created_at ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []*models.Investment
	for rows.Next() {
		var investment models.Investment
		err := rows.Scan(
			&investment.ID,
			&investment.LoanID,
			&investment.InvestorID,
			&investment.Amount,
			&investment.InvestmentDate,
			&investment.ExpectedReturn,
			&investment.AgreementSent,
			&investment.AgreementSentAt,
			&investment.CreatedAt,
			&investment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		investments = append(investments, &investment)
	}

	return investments, nil
}

// UpdateInvestmentAgreementSent updates the agreement sent status for an investment
func (r *LoanRepository) UpdateInvestmentAgreementSent(investmentID uuid.UUID, agreementSent bool, agreementSentAt *time.Time) error {
	query := `UPDATE investments SET agreement_sent = $1, agreement_sent_at = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`

	_, err := r.db.Exec(query, agreementSent, agreementSentAt, investmentID)
	return err
}

// GetInvestorByID gets an investor by ID
func (r *LoanRepository) GetInvestorByID(investorID uuid.UUID) (*models.Investor, error) {
	query := `SELECT id, investor_code, name, email, phone_number, is_active, created_at, updated_at
			  FROM investors WHERE id = $1 AND deleted_at IS NULL`

	var investor models.Investor
	err := r.db.QueryRow(query, investorID).Scan(
		&investor.ID,
		&investor.InvestorCode,
		&investor.Name,
		&investor.Email,
		&investor.PhoneNumber,
		&investor.IsActive,
		&investor.CreatedAt,
		&investor.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &investor, nil
}

// GetBorrowerByID gets a borrower by ID
func (r *LoanRepository) GetBorrowerByID(borrowerID uuid.UUID) (*models.Borrower, error) {
	query := `SELECT id, id_number, first_name, last_name, email, phone_number, address, created_at, updated_at
			  FROM borrowers WHERE id = $1 AND deleted_at IS NULL`

	var borrower models.Borrower
	err := r.db.QueryRow(query, borrowerID).Scan(
		&borrower.ID,
		&borrower.IDNumber,
		&borrower.FirstName,
		&borrower.LastName,
		&borrower.Email,
		&borrower.PhoneNumber,
		&borrower.Address,
		&borrower.CreatedAt,
		&borrower.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &borrower, nil
}

// CreateEmailNotification creates an email notification record
func (r *LoanRepository) CreateEmailNotification(notification *models.EmailNotification) error {
	query := `INSERT INTO email_notifications (id, investor_id, loan_id, email_type, email_subject, email_body, sent_at, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING created_at, updated_at`

	return r.db.QueryRow(query,
		notification.ID,
		notification.InvestorID,
		notification.LoanID,
		notification.EmailType,
		notification.EmailSubject,
		notification.EmailBody,
		notification.SentAt,
		notification.Status,
	).Scan(&notification.CreatedAt, &notification.UpdatedAt)
}

// UpdateLoanAgreementLetterURL updates the agreement letter URL for a loan
func (r *LoanRepository) UpdateLoanAgreementLetterURL(tx *sql.Tx, loanID uuid.UUID, agreementURL string) error {
	query := `UPDATE loans SET agreement_letter_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, agreementURL, loanID)
	} else {
		_, err = r.db.Exec(query, agreementURL, loanID)
	}
	return err
}
