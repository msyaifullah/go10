// internal/services/disbursement_service.go
package services

import (
	"database/sql"
	"fmt"
	"time"

	"loan-service/internal/constant"
	"loan-service/internal/models"
	"loan-service/internal/repositories"
	"loan-service/pkg/adapters"
	"loan-service/pkg/logger"

	"github.com/google/uuid"
)

type LoanService struct {
	loanRepo       repositories.LoanRepositoryInterface
	paymentAdapter adapters.PaymentAdapterInterface
	emailAdapter   adapters.EmailAdapterInterface
	logger         logger.LoggerInterface
	db             *sql.DB
}

func NewLoanService(
	loanRepo repositories.LoanRepositoryInterface,
	paymentAdapter adapters.PaymentAdapterInterface,
	emailAdapter adapters.EmailAdapterInterface,
	logger logger.LoggerInterface,
	db *sql.DB,
) LoanServiceInterface {
	return &LoanService{
		loanRepo:       loanRepo,
		paymentAdapter: paymentAdapter,
		emailAdapter:   emailAdapter,
		logger:         logger,
		db:             db,
	}
}

func (s *LoanService) withTransaction(fn func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Error("Failed to begin transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			s.logger.Error("Panic in transaction, rolling back", map[string]interface{}{
				"panic": p,
			})
			tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		s.logger.Error("Transaction failed, rolling back", map[string]interface{}{
			"error": err.Error(),
		})
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			s.logger.Error("Failed to rollback transaction", map[string]interface{}{
				"rollback_error": rollbackErr.Error(),
				"original_error": err.Error(),
			})
			return fmt.Errorf("failed to rollback transaction: %w (original error: %w)", rollbackErr, err)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *LoanService) GetLoanByID(id uuid.UUID) (*models.LoanSummaryResponse, error) {
	s.logger.Info("Getting loan by ID", map[string]interface{}{"id": id})

	loan, err := s.loanRepo.GetLoanByID(nil, id)
	if err != nil {
		s.logger.Error("Failed to get loan by ID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return &models.LoanSummaryResponse{
		ID:                  loan.ID,
		BorrowerName:        loan.Borrower.FullName(),
		PrincipalAmount:     loan.PrincipalAmount,
		InterestRate:        loan.InterestRate,
		ROI:                 loan.ROI,
		State:               loan.State,
		TotalInvested:       loan.TotalInvested,
		RemainingInvestment: loan.RemainingInvestmentAmount(),
		InvestorCount:       len(loan.Investments),
		CreatedAt:           loan.CreatedAt,
		UpdatedAt:           loan.UpdatedAt,
	}, nil
}

func (s *LoanService) ProcessCreateLoan(req *models.CreateLoanRequest) (*models.LoanSummaryResponse, error) {
	s.logger.Info("Processing loan creation", map[string]interface{}{"request": req})

	loan := &models.Loan{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		BorrowerID:      req.BorrowerID,
		PrincipalAmount: req.PrincipalAmount,
		InterestRate:    req.InterestRate,
		ROI:             req.ROI,
		State:           models.LoanStateProposed,
		TotalInvested:   0,
	}

	loan, err := s.loanRepo.CreateLoan(nil, loan)
	if err != nil {
		s.logger.Error("Failed to create loan", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info("Loan created", map[string]interface{}{"loan": loan})

	// Fetch the loan with borrower data to build the response
	loanWithBorrower, err := s.loanRepo.GetLoanByID(nil, loan.ID)
	if err != nil {
		s.logger.Error("Failed to fetch loan with borrower data", map[string]interface{}{
			"error":   err.Error(),
			"loan_id": loan.ID.String(),
		})
		return nil, err
	}

	return &models.LoanSummaryResponse{
		ID:                  loanWithBorrower.ID,
		BorrowerName:        loanWithBorrower.Borrower.FullName(),
		PrincipalAmount:     loanWithBorrower.PrincipalAmount,
		InterestRate:        loanWithBorrower.InterestRate,
		ROI:                 loanWithBorrower.ROI,
		State:               loanWithBorrower.State,
		TotalInvested:       loanWithBorrower.TotalInvested,
		RemainingInvestment: loanWithBorrower.RemainingInvestmentAmount(),
		InvestorCount:       len(loanWithBorrower.Investments),
		CreatedAt:           loanWithBorrower.CreatedAt,
		UpdatedAt:           loanWithBorrower.UpdatedAt,
	}, nil
}

func (s *LoanService) ProcessApproveLoan(id uuid.UUID, req *models.CreateApprovalRequest) (*models.LoanApprovalResponse, error) {
	s.logger.Info("Approving loan", map[string]interface{}{"id": id, "request": req})

	// Use transaction to ensure data consistency
	var result *models.LoanApprovalResponse
	err := s.withTransaction(func(tx *sql.Tx) error {
		var approvalErr error
		result, approvalErr = s.processApproveLoanTx(tx, id, req)
		return approvalErr
	})

	return result, err
}

func (s *LoanService) processApproveLoanTx(tx *sql.Tx, id uuid.UUID, req *models.CreateApprovalRequest) (*models.LoanApprovalResponse, error) {
	approval := &models.Approval{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		LoanID:              id,
		ValidatorID:         req.ValidatorID,
		ApprovalDate:        req.ApprovalDate,
		VisitProofImageURL:  req.VisitProofImageURL,
		VisitProofImageType: models.FileType(req.VisitProofImageType),
		Notes:               req.Notes,
	}

	// check if loan is already approved
	loan, err := s.loanRepo.GetLoanByID(tx, id)
	if err != nil {
		s.logger.Error("Failed to get loan by ID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Validate state transition with detailed error messages
	if err := loan.ValidateStateTransition(models.LoanStateApproved); err != nil {
		s.logger.Error("State transition validation failed", map[string]interface{}{
			"error":         err.Error(),
			"loan_id":       id.String(),
			"current_state": loan.State.String(),
			"target_state":  models.LoanStateApproved.String(),
		})
		return nil, fmt.Errorf("loan approval validation failed: %w", err)
	}

	approval, err = s.loanRepo.CreateApproval(tx, approval)
	if err != nil {
		s.logger.Error("Failed to create approval", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info("Approval created", map[string]interface{}{"approval": approval})

	state, err := s.loanRepo.UpdateLoanState(tx, approval.LoanID, models.LoanStateApproved)
	if err != nil {
		s.logger.Error("Failed to update loan state", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info("Loan state updated", map[string]interface{}{"state": state})

	record, err := s.loanRepo.RecordLoanStateHistory(tx, models.LoanStateProposed, state, approval.ValidatorID, "Loan approved")
	if err != nil {
		s.logger.Error("Failed to record loan state history", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info("Loan state history recorded", map[string]interface{}{"record": record})

	return &models.LoanApprovalResponse{
		ID:                  approval.ID,
		LoanID:              approval.LoanID,
		ValidatorID:         approval.ValidatorID,
		ApprovalDate:        approval.ApprovalDate,
		VisitProofImageURL:  approval.VisitProofImageURL,
		VisitProofImageType: approval.VisitProofImageType,
		Notes:               approval.Notes,
		CreatedAt:           approval.CreatedAt,
		UpdatedAt:           approval.UpdatedAt,
	}, nil
}

func (s *LoanService) ProcessInvestment(loanID uuid.UUID, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error) {
	s.logger.Info("Processing investment", map[string]interface{}{"loan_id": loanID, "request": req})

	// Use transaction to ensure data consistency
	var result *models.InvestmentResponse
	err := s.withTransaction(func(tx *sql.Tx) error {
		var investmentErr error
		result, investmentErr = s.processInvestmentTx(tx, loanID, req)
		return investmentErr
	})

	return result, err
}

func (s *LoanService) processInvestmentTx(tx *sql.Tx, loanID uuid.UUID, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error) {
	// Get the loan with current state for validation
	loan, err := s.loanRepo.GetLoanByID(tx, loanID)
	if err != nil {
		s.logger.Error("Failed to get loan by ID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Use the existing ValidateInvestmentAmount for business rule validation
	if err := loan.ValidateInvestmentAmount(req.Amount); err != nil {
		s.logger.Error("Investment amount validation failed", map[string]interface{}{
			"error":   err.Error(),
			"loan_id": loanID.String(),
			"amount":  req.Amount,
			"state":   loan.State.String(),
		})
		return nil, fmt.Errorf("investment validation failed: %w", err)
	}

	// Create investment record first
	investment := &models.Investment{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		LoanID:         loanID,
		InvestorID:     req.InvestorID,
		Amount:         req.Amount,
		ExpectedReturn: req.Amount * (1 + loan.ROI), // Calculate expected return based on ROI
		InvestmentDate: req.InvestmentDate,
	}

	investment, err = s.loanRepo.CreateInvestment(tx, investment)
	if err != nil {
		s.logger.Error("Failed to create investment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Use atomic update to prevent race conditions
	// This is the key fix - combines validation and update in one atomic operation
	updatedLoan, err := s.loanRepo.UpdateLoanTotalInvested(tx, loanID, req.Amount)
	if err != nil {
		s.logger.Error("Failed to update loan total invested atomically", map[string]interface{}{
			"error":   err.Error(),
			"loan_id": loanID.String(),
			"amount":  req.Amount,
		})
		return nil, fmt.Errorf("investment processing failed: %w", err)
	}

	// Check if loan is now fully invested and can transition to invested state
	if updatedLoan.IsFullyInvested() && updatedLoan.State == models.LoanStateApproved {
		// Validate state transition to invested
		if err := updatedLoan.ValidateStateTransition(models.LoanStateInvested); err != nil {
			s.logger.Error("Cannot transition to invested state", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, fmt.Errorf("cannot transition to invested state: %w", err)
		}

		// Update loan state to invested
		newState, err := s.loanRepo.UpdateLoanState(tx, loanID, models.LoanStateInvested)
		if err != nil {
			s.logger.Error("Failed to update loan state to invested", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		// Update loan agreement letter URL
		agreementURL, err := s.mockAgreementLetterURL(newState)
		if err != nil {
			s.logger.Error("Failed to generate loan agreement letter URL", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		// Update loan agreement letter URL
		err = s.loanRepo.UpdateLoanAgreementLetterURL(tx, loanID, agreementURL)
		if err != nil {
			s.logger.Error("Failed to update loan agreement letter URL", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		// Record state history
		_, err = s.loanRepo.RecordLoanStateHistory(tx, models.LoanStateApproved, newState, uuid.MustParse(constant.SystemEmployeeID), "Investment target achieved")
		if err != nil {
			s.logger.Error("Failed to record loan state history", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		s.logger.Info("Loan state updated to invested", map[string]interface{}{
			"loan_id":   loanID.String(),
			"new_state": newState.State.String(),
		})
	}

	return &models.InvestmentResponse{
		ID:             investment.ID,
		LoanID:         investment.LoanID,
		InvestorID:     investment.InvestorID,
		Amount:         investment.Amount,
		ExpectedReturn: investment.ExpectedReturn,
		InvestmentDate: investment.InvestmentDate,
		CreatedAt:      investment.CreatedAt,
		UpdatedAt:      investment.UpdatedAt,
	}, nil
}

func (s *LoanService) ProcessDisbursement(loanID uuid.UUID, req *models.CreateDisbursementRequest) (*models.DisbursementResponse, error) {
	s.logger.Info("Processing disbursement", map[string]interface{}{"loan_id": loanID, "request": req})

	// Use transaction to ensure data consistency
	var result *models.DisbursementResponse
	err := s.withTransaction(func(tx *sql.Tx) error {
		var disbursementErr error
		result, disbursementErr = s.processDisbursementTx(tx, loanID, req)
		return disbursementErr
	})

	return result, err
}

func (s *LoanService) processDisbursementTx(tx *sql.Tx, loanID uuid.UUID, req *models.CreateDisbursementRequest) (*models.DisbursementResponse, error) {
	// Get the loan with current state
	loan, err := s.loanRepo.GetLoanByID(tx, loanID)
	if err != nil {
		s.logger.Error("Failed to get loan by ID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Validate state transition to disbursed
	if err := loan.ValidateStateTransition(models.LoanStateDisbursed); err != nil {
		s.logger.Error("State transition validation failed", map[string]interface{}{
			"error":         err.Error(),
			"loan_id":       loanID.String(),
			"current_state": loan.State.String(),
			"target_state":  models.LoanStateDisbursed.String(),
		})
		return nil, fmt.Errorf("disbursement validation failed: %w", err)
	}

	// Validate disbursed amount matches principal amount
	if req.DisbursedAmount != loan.PrincipalAmount {
		return nil, fmt.Errorf("disbursed amount %.2f must match principal amount %.2f",
			req.DisbursedAmount, loan.PrincipalAmount)
	}

	// Create disbursement record
	disbursement := &models.Disbursement{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		LoanID:                  loanID,
		FieldOfficerID:          req.FieldOfficerID,
		DisbursementDate:        req.DisbursementDate,
		SignedAgreementURL:      req.SignedAgreementURL,
		SignedAgreementFileType: models.FileType(req.SignedAgreementFileType),
		DisbursedAmount:         req.DisbursedAmount,
		Notes:                   req.Notes,
	}

	disbursement, err = s.loanRepo.CreateDisbursement(tx, disbursement)
	if err != nil {
		s.logger.Error("Failed to create disbursement", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Update loan state to disbursed
	newState, err := s.loanRepo.UpdateLoanState(tx, loanID, models.LoanStateDisbursed)
	if err != nil {
		s.logger.Error("Failed to update loan state to disbursed", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Record state history
	_, err = s.loanRepo.RecordLoanStateHistory(tx, models.LoanStateInvested, newState, req.FieldOfficerID, "Loan disbursed")
	if err != nil {
		s.logger.Error("Failed to record loan state history", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.logger.Info("Loan disbursed on process disbursement", map[string]interface{}{
		"loan_id":   loanID.String(),
		"new_state": newState.State.String(),
	})

	// Process payment for the disbursed amount
	paymentResult, err := s.paymentAdapter.ProcessPayment(req.DisbursedAmount, "disbursement_token_"+loanID.String())
	if err != nil {
		s.logger.Error("Failed to process payment for disbursement", map[string]interface{}{
			"error":            err.Error(),
			"loan_id":          loanID.String(),
			"disbursed_amount": req.DisbursedAmount,
		})
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	s.logger.Info("Payment processed successfully for disbursement", map[string]interface{}{
		"loan_id":        loanID.String(),
		"transaction_id": paymentResult.TransactionID,
		"status":         paymentResult.Status,
		"amount":         req.DisbursedAmount,
	})

	return &models.DisbursementResponse{
		ID:                      disbursement.ID,
		LoanID:                  disbursement.LoanID,
		FieldOfficerID:          disbursement.FieldOfficerID,
		DisbursementDate:        disbursement.DisbursementDate,
		SignedAgreementURL:      disbursement.SignedAgreementURL,
		SignedAgreementFileType: disbursement.SignedAgreementFileType,
		DisbursedAmount:         disbursement.DisbursedAmount,
		Notes:                   disbursement.Notes,
		CreatedAt:               disbursement.CreatedAt,
		UpdatedAt:               disbursement.UpdatedAt,
	}, nil
}

// mockAgreementLetterURL generates a URL for the loan agreement letter
func (s *LoanService) mockAgreementLetterURL(loan *models.Loan) (string, error) {
	// Generate a unique filename for the agreement letter
	filename := fmt.Sprintf("agreement_loan_%s.pdf", loan.ID.String()[:8])

	// Construct the URL - this could be a template URL or generated based on configuration
	// For now, using a simple template URL structure
	agreementURL := fmt.Sprintf("https://storage.example.com/agreements?loan_id=%s&filename=%s&token=%s", loan.ID.String(), filename, "mock-token-123")

	s.logger.Info("Generated agreement letter URL", map[string]interface{}{
		"loan_id":       loan.ID.String(),
		"filename":      filename,
		"agreement_url": agreementURL,
	})

	return agreementURL, nil
}
