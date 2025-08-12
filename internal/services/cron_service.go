package services

import (
	"database/sql"
	"fmt"
	"time"

	"loan-service/internal/models"
	"loan-service/internal/repositories"
	"loan-service/pkg/adapters"
	"loan-service/pkg/config"
	"loan-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	loanRepo     repositories.LoanRepositoryInterface
	emailAdapter adapters.EmailAdapterInterface
	logger       *logger.Logger
	db           *sql.DB
	cron         *cron.Cron
	config       *config.Config
}

func NewCronService(
	loanRepo repositories.LoanRepositoryInterface,
	emailAdapter adapters.EmailAdapterInterface,
	logger *logger.Logger,
	db *sql.DB,
	config *config.Config,
) *CronService {
	return &CronService{
		loanRepo:     loanRepo,
		emailAdapter: emailAdapter,
		logger:       logger,
		db:           db,
		cron:         cron.New(cron.WithSeconds()),
		config:       config,
	}
}

// Start starts the cron service and schedules jobs
func (s *CronService) Start() {
	s.logger.Info("Starting cron service", map[string]interface{}{})

	// Schedule investment agreement email job using configuration
	schedule := s.config.Cron.InvestmentAgreementSchedule
	if schedule == "" {
		schedule = "0 */5 * * * *" // Default fallback
		s.logger.Warn("Using default cron schedule for investment agreements", map[string]interface{}{
			"schedule": schedule,
		})
	}

	_, err := s.cron.AddFunc(schedule, s.processInvestmentAgreements)
	if err != nil {
		s.logger.Error("Failed to schedule investment agreement job", map[string]interface{}{
			"error":    err.Error(),
			"schedule": schedule,
		})
		return
	}

	s.cron.Start()
	s.logger.Info("Cron service started successfully", map[string]interface{}{
		"investment_agreement_schedule": schedule,
	})
}

// Stop stops the cron service
func (s *CronService) Stop() {
	s.logger.Info("Stopping cron service", map[string]interface{}{})
	s.cron.Stop()
}

// processInvestmentAgreements processes investments that need agreement emails sent
func (s *CronService) processInvestmentAgreements() {
	s.logger.Info("Starting investment agreement processing job", map[string]interface{}{})

	// Get investments that need agreement emails
	investments, err := s.loanRepo.GetInvestmentsNeedingAgreementEmail()
	if err != nil {
		s.logger.Error("Failed to get investments needing agreement emails", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(investments) == 0 {
		s.logger.Info("No investments need agreement emails", map[string]interface{}{})
		return
	}

	s.logger.Info("Found investments needing agreement emails", map[string]interface{}{
		"count": len(investments),
	})

	// Group investments by loan ID to optimize database calls
	investmentsByLoan := make(map[uuid.UUID][]*models.Investment)
	for _, investment := range investments {
		investmentsByLoan[investment.LoanID] = append(investmentsByLoan[investment.LoanID], investment)
	}

	s.logger.Info("Grouped investments by loan", map[string]interface{}{
		"unique_loans": len(investmentsByLoan),
	})

	// Process investments grouped by loan
	for loanID, loanInvestments := range investmentsByLoan {
		if err := s.processInvestmentsForLoan(loanID, loanInvestments); err != nil {
			s.logger.Error("Failed to process investments for loan", map[string]interface{}{
				"loan_id": loanID.String(),
				"error":   err.Error(),
			})
			continue
		}
	}

	s.logger.Info("Investment agreement processing job completed", map[string]interface{}{
		"processed_count": len(investments),
	})
}

// processInvestmentsForLoan processes all investments for a single loan
func (s *CronService) processInvestmentsForLoan(loanID uuid.UUID, investments []*models.Investment) error {
	s.logger.Info("Processing investments for loan", map[string]interface{}{
		"loan_id":          loanID.String(),
		"investment_count": len(investments),
	})

	// Get loan details once for all investments
	loan, err := s.loanRepo.GetLoanByID(nil, loanID)
	if err != nil {
		return fmt.Errorf("failed to get loan: %w", err)
	}

	// Get borrower details once for all investments
	borrower, err := s.loanRepo.GetBorrowerByID(loan.BorrowerID)
	if err != nil {
		return fmt.Errorf("failed to get borrower: %w", err)
	}

	// Process each investment for this loan
	for _, investment := range investments {
		if err := s.processSendingEmail(investment, loan, borrower); err != nil {
			s.logger.Error("Failed to process sending email", map[string]interface{}{
				"investment_id": investment.ID.String(),
				"loan_id":       loanID.String(),
				"error":         err.Error(),
			})
			continue
		}
	}

	return nil
}

// processSendingEmail processes sending email for a single investment with pre-fetched loan and borrower data
func (s *CronService) processSendingEmail(investment *models.Investment, loan *models.Loan, borrower *models.Borrower) error {
	s.logger.Info("Processing investment agreement", map[string]interface{}{
		"investment_id": investment.ID.String(),
		"loan_id":       investment.LoanID.String(),
		"investor_id":   investment.InvestorID.String(),
	})

	// Get investor details
	investor, err := s.loanRepo.GetInvestorByID(investment.InvestorID)
	if err != nil {
		return fmt.Errorf("failed to get investor: %w", err)
	}

	// Prepare email content
	subject := fmt.Sprintf("Investment Agreement - Loan #%s", loan.ID.String()[:8])
	body := s.emailAdapter.GenerateAgreementEmailBody(investment, loan, investor, borrower)

	// Send email
	if err := s.emailAdapter.SendEmail(investor.Email, subject, body); err != nil {
		s.logger.Error("Failed to send agreement email", map[string]interface{}{
			"investment_id":  investment.ID.String(),
			"investor_email": investor.Email,
			"error":          err.Error(),
		})
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Update investment to mark agreement as sent
	now := time.Now()
	if err := s.loanRepo.UpdateInvestmentAgreementSent(investment.ID, true, &now); err != nil {
		s.logger.Error("Failed to update investment agreement sent status", map[string]interface{}{
			"investment_id": investment.ID.String(),
			"error":         err.Error(),
		})
		return fmt.Errorf("failed to update investment: %w", err)
	}

	// Record email notification
	notification := &models.EmailNotification{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		InvestorID:   investment.InvestorID,
		LoanID:       investment.LoanID,
		EmailType:    "agreement_notification",
		EmailSubject: subject,
		EmailBody:    body,
		SentAt:       now,
		Status:       "sent",
	}

	if err := s.loanRepo.CreateEmailNotification(notification); err != nil {
		s.logger.Error("Failed to create email notification record", map[string]interface{}{
			"investment_id": investment.ID.String(),
			"error":         err.Error(),
		})
		// Don't fail the entire process for notification recording failure
	}

	s.logger.Info("Investment agreement processed successfully", map[string]interface{}{
		"investment_id":  investment.ID.String(),
		"investor_email": investor.Email,
	})

	return nil
}
