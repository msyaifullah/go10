package services

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"loan-service/internal/models"
	"loan-service/pkg/adapters"

	"sync"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) CreateLoan(tx *sql.Tx, loan *models.Loan) (*models.Loan, error) {
	args := m.Called(tx, loan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetLoanByID(tx *sql.Tx, loanID uuid.UUID) (*models.Loan, error) {
	args := m.Called(tx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) UpdateLoanState(tx *sql.Tx, loanID uuid.UUID, newState models.LoanState) (*models.Loan, error) {
	args := m.Called(tx, loanID, newState)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) RecordLoanStateHistory(tx *sql.Tx, prevState models.LoanState, loan *models.Loan, employeeID uuid.UUID, changeReason string) (*models.LoanStateHistory, error) {
	args := m.Called(tx, prevState, loan, employeeID, changeReason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanStateHistory), args.Error(1)
}

func (m *MockLoanRepository) CreateApproval(tx *sql.Tx, approval *models.Approval) (*models.Approval, error) {
	args := m.Called(tx, approval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Approval), args.Error(1)
}

func (m *MockLoanRepository) CreateInvestment(tx *sql.Tx, investment *models.Investment) (*models.Investment, error) {
	args := m.Called(tx, investment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Investment), args.Error(1)
}

func (m *MockLoanRepository) UpdateLoanTotalInvested(tx *sql.Tx, loanID uuid.UUID, newTotal float64) (*models.Loan, error) {
	args := m.Called(tx, loanID, newTotal)
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetInvestmentsNeedingAgreementEmail() ([]*models.Investment, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Investment), args.Error(1)
}

func (m *MockLoanRepository) UpdateInvestmentAgreementSent(investmentID uuid.UUID, agreementSent bool, agreementSentAt *time.Time) error {
	args := m.Called(investmentID, agreementSent, agreementSentAt)
	return args.Error(0)
}

func (m *MockLoanRepository) UpdateLoanAgreementLetterURL(tx *sql.Tx, loanID uuid.UUID, agreementURL string) error {
	args := m.Called(tx, loanID, agreementURL)
	return args.Error(0)
}

func (m *MockLoanRepository) CreateDisbursement(tx *sql.Tx, disbursement *models.Disbursement) (*models.Disbursement, error) {
	args := m.Called(tx, disbursement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Disbursement), args.Error(1)
}

func (m *MockLoanRepository) GetInvestorByID(investorID uuid.UUID) (*models.Investor, error) {
	args := m.Called(investorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Investor), args.Error(1)
}

func (m *MockLoanRepository) GetBorrowerByID(borrowerID uuid.UUID) (*models.Borrower, error) {
	args := m.Called(borrowerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Borrower), args.Error(1)
}

func (m *MockLoanRepository) CreateEmailNotification(notification *models.EmailNotification) error {
	args := m.Called(notification)
	return args.Error(0)
}

type MockPaymentAdapter struct {
	mock.Mock
}

func (m *MockPaymentAdapter) ProcessPayment(amount float64, token string) (*adapters.PaymentResult, error) {
	args := m.Called(amount, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*adapters.PaymentResult), args.Error(1)
}

type MockEmailAdapter struct {
	mock.Mock
}

func (m *MockEmailAdapter) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *MockEmailAdapter) GenerateAgreementEmailBody(
	investment *models.Investment,
	loan *models.Loan,
	investor *models.Investor,
	borrower *models.Borrower,
) string {
	args := m.Called(investment, loan, investor, borrower)
	return args.String(0)
}

// SilentLogger is a logger that does nothing - perfect for tests
type TestLogger struct{}

func (s *TestLogger) Info(message string, data map[string]interface{})  {}
func (s *TestLogger) Error(message string, data map[string]interface{}) {}
func (s *TestLogger) Debug(message string, data map[string]interface{}) {}
func (s *TestLogger) Warn(message string, data map[string]interface{})  {}
func (s *TestLogger) Fatal(message string, data map[string]interface{}) {}

// TestLoanService is a test-specific version that overrides withTransaction
type TestLoanService struct {
	*LoanService
}

// Override withTransaction to bypass actual transactions in tests
func (s *TestLoanService) withTransaction(fn func(*sql.Tx) error) error {
	return fn(nil) // Pass nil transaction to simulate no transaction
}

// Override methods that use transactions to use our mocked withTransaction
func (s *TestLoanService) ProcessApproveLoan(id uuid.UUID, req *models.CreateApprovalRequest) (*models.LoanApprovalResponse, error) {
	// Use transaction to ensure data consistency
	var result *models.LoanApprovalResponse
	err := s.withTransaction(func(tx *sql.Tx) error {
		var approvalErr error
		result, approvalErr = s.processApproveLoanTx(tx, id, req)
		return approvalErr
	})

	return result, err
}

func (s *TestLoanService) ProcessInvestment(loanID uuid.UUID, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error) {
	// Use transaction to ensure data consistency
	var result *models.InvestmentResponse
	err := s.withTransaction(func(tx *sql.Tx) error {
		var investmentErr error
		result, investmentErr = s.processInvestmentTx(tx, loanID, req)
		return investmentErr
	})

	return result, err
}

func (s *TestLoanService) ProcessDisbursement(loanID uuid.UUID, req *models.CreateDisbursementRequest) (*models.DisbursementResponse, error) {
	// Use transaction to ensure data consistency
	var result *models.DisbursementResponse
	err := s.withTransaction(func(tx *sql.Tx) error {
		var disbursementErr error
		result, disbursementErr = s.processDisbursementTx(tx, loanID, req)
		return disbursementErr
	})

	return result, err
}

// Test setup helper - now uses mocked dependencies with silent logger
func setupTestLoanService() (*TestLoanService, *MockLoanRepository, *MockPaymentAdapter, *MockEmailAdapter) {
	mockRepo := &MockLoanRepository{}
	mockPayment := &MockPaymentAdapter{}
	mockEmail := &MockEmailAdapter{}

	// Use silent logger to eliminate log messages
	silentLogger := &TestLogger{}

	// Create a minimal database connection for testing
	// We'll use a nil DB since the service will use transactions that we mock
	var db *sql.DB

	// Create the real LoanService with mocked dependencies
	baseService := NewLoanService(mockRepo, mockPayment, mockEmail, silentLogger, db).(*LoanService)

	// Wrap it in TestLoanService to override withTransaction
	service := &TestLoanService{LoanService: baseService}

	return service, mockRepo, mockPayment, mockEmail
}

// Helper function to create test data
func createTestLoan(id uuid.UUID, state models.LoanState, totalInvested float64) *models.Loan {
	return &models.Loan{
		BaseModel: models.BaseModel{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		BorrowerID:      uuid.New(),
		PrincipalAmount: 10000.0,
		InterestRate:    0.10,
		ROI:             0.15,
		State:           state,
		TotalInvested:   totalInvested,
		Borrower: &models.Borrower{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},
			FirstName: "John",
			LastName:  "Doe",
		},
		Investments: []models.Investment{},
	}
}

// ===== SIMPLIFIED TESTS - SUCCESS CASES ONLY =====

func TestLoanService_GetLoanByID_Success(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	loan := createTestLoan(loanID, models.LoanStateProposed, 0)

	mockRepo.On("GetLoanByID", (*sql.Tx)(nil), loanID).Return(loan, nil)

	result, err := service.GetLoanByID(loanID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, loanID, result.ID)
	assert.Equal(t, "John Doe", result.BorrowerName)
	assert.Equal(t, 10000.0, result.PrincipalAmount)
	assert.Equal(t, models.LoanStateProposed, result.State)

	mockRepo.AssertExpectations(t)
}

func TestLoanService_GetLoanByID_Error(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	mockRepo.On("GetLoanByID", (*sql.Tx)(nil), loanID).Return(nil, errors.New("loan not found"))

	result, err := service.GetLoanByID(loanID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "loan not found")

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessCreateLoan_Success(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	req := &models.CreateLoanRequest{
		BorrowerID:      uuid.New(),
		PrincipalAmount: 10000.0,
		InterestRate:    0.10,
		ROI:             0.15,
	}

	createdLoan := createTestLoan(uuid.New(), models.LoanStateProposed, 0)
	loanWithBorrower := createTestLoan(createdLoan.ID, models.LoanStateProposed, 0)

	mockRepo.On("CreateLoan", (*sql.Tx)(nil), mock.AnythingOfType("*models.Loan")).Return(createdLoan, nil)
	mockRepo.On("GetLoanByID", (*sql.Tx)(nil), createdLoan.ID).Return(loanWithBorrower, nil)

	result, err := service.ProcessCreateLoan(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.PrincipalAmount, result.PrincipalAmount)
	assert.Equal(t, req.InterestRate, result.InterestRate)
	assert.Equal(t, req.ROI, result.ROI)
	assert.Equal(t, models.LoanStateProposed, result.State)

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessCreateLoan_Error(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	req := &models.CreateLoanRequest{
		BorrowerID:      uuid.New(),
		PrincipalAmount: 10000.0,
		InterestRate:    0.10,
		ROI:             0.15,
	}

	mockRepo.On("CreateLoan", (*sql.Tx)(nil), mock.AnythingOfType("*models.Loan")).Return(nil, errors.New("database error"))

	result, err := service.ProcessCreateLoan(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessApproveLoan_Success(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	req := &models.CreateApprovalRequest{
		ValidatorID:         uuid.New(),
		ApprovalDate:        time.Now(),
		VisitProofImageURL:  "https://example.com/proof.jpg",
		VisitProofImageType: models.FileTypeJPEG,
		Notes:               "Approved after site visit",
	}

	loan := createTestLoan(loanID, models.LoanStateProposed, 0)
	approval := &models.Approval{
		BaseModel:           models.BaseModel{ID: uuid.New()},
		LoanID:              loanID,
		ValidatorID:         req.ValidatorID,
		ApprovalDate:        req.ApprovalDate,
		VisitProofImageURL:  req.VisitProofImageURL,
		VisitProofImageType: req.VisitProofImageType,
		Notes:               req.Notes,
	}
	updatedLoan := createTestLoan(loanID, models.LoanStateApproved, 0)
	history := &models.LoanStateHistory{
		BaseModel: models.BaseModel{ID: uuid.New()},
		LoanID:    loanID,
	}

	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(loan, nil)
	mockRepo.On("CreateApproval", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("*models.Approval")).Return(approval, nil)
	mockRepo.On("UpdateLoanState", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("uuid.UUID"), models.LoanStateApproved).Return(updatedLoan, nil)
	mockRepo.On("RecordLoanStateHistory", mock.AnythingOfType("*sql.Tx"), models.LoanStateProposed, mock.AnythingOfType("*models.Loan"), mock.AnythingOfType("uuid.UUID"), "Loan approved").Return(history, nil)

	result, err := service.ProcessApproveLoan(loanID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, loanID, result.LoanID)
	assert.NotEqual(t, uuid.Nil, result.ValidatorID)
	assert.NotEmpty(t, result.VisitProofImageURL)

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessApproveLoan_Error(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	req := &models.CreateApprovalRequest{
		ValidatorID:         uuid.New(),
		ApprovalDate:        time.Now(),
		VisitProofImageURL:  "https://example.com/proof.jpg",
		VisitProofImageType: models.FileTypeJPEG,
		Notes:               "Approved after site visit",
	}

	// Loan already approved - should fail
	loan := createTestLoan(loanID, models.LoanStateApproved, 0)

	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(loan, nil)

	result, err := service.ProcessApproveLoan(loanID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot transition")

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessInvestment_Success(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	req := &models.CreateInvestmentRequest{
		InvestorID:     uuid.New(),
		Amount:         5000.0,
		InvestmentDate: time.Now(),
	}

	loan := createTestLoan(loanID, models.LoanStateApproved, 0)
	investment := &models.Investment{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		LoanID:         loanID,
		InvestorID:     req.InvestorID,
		Amount:         req.Amount,
		ExpectedReturn: req.Amount * 1.15, // 15% ROI
		InvestmentDate: req.InvestmentDate,
	}
	updatedLoan := createTestLoan(loanID, models.LoanStateApproved, 5000.0)

	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(loan, nil)
	mockRepo.On("CreateInvestment", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("*models.Investment")).Return(investment, nil)
	mockRepo.On("UpdateLoanTotalInvested", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("uuid.UUID"), 5000.0).Return(updatedLoan, nil)

	result, err := service.ProcessInvestment(loanID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, loanID, result.LoanID)
	assert.Equal(t, req.Amount, result.Amount)
	assert.NotEqual(t, uuid.Nil, result.InvestorID)

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessInvestment_Error(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	req := &models.CreateInvestmentRequest{
		InvestorID:     uuid.New(),
		Amount:         15000.0, // More than principal amount
		InvestmentDate: time.Now(),
	}

	loan := createTestLoan(loanID, models.LoanStateApproved, 0)

	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(loan, nil)

	result, err := service.ProcessInvestment(loanID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "exceeds remaining amount")

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessInvestment_RaceConditionHandling(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	req1 := &models.CreateInvestmentRequest{
		InvestorID:     uuid.New(),
		Amount:         3000.0,
		InvestmentDate: time.Now(),
	}
	req2 := &models.CreateInvestmentRequest{
		InvestorID:     uuid.New(),
		Amount:         2000.0,
		InvestmentDate: time.Now(),
	}

	// Initial loan state
	initialLoan := createTestLoan(loanID, models.LoanStateApproved, 0)

	// Investment records
	investment1 := &models.Investment{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		LoanID:         loanID,
		InvestorID:     req1.InvestorID,
		Amount:         req1.Amount,
		ExpectedReturn: req1.Amount * 1.15,
		InvestmentDate: req1.InvestmentDate,
	}
	investment2 := &models.Investment{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		LoanID:         loanID,
		InvestorID:     req2.InvestorID,
		Amount:         req2.Amount,
		ExpectedReturn: req2.Amount * 1.15,
		InvestmentDate: req2.InvestmentDate,
	}

	// Updated loan states after atomic updates
	updatedLoan1 := createTestLoan(loanID, models.LoanStateApproved, 3000.0)
	updatedLoan2 := createTestLoan(loanID, models.LoanStateApproved, 5000.0)

	// Final loan state when fully invested
	investedLoan := createTestLoan(loanID, models.LoanStateInvested, 5000.0)

	// Mock expectations for first investment
	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(initialLoan, nil).Once()
	mockRepo.On("CreateInvestment", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("*models.Investment")).Return(investment1, nil).Once()
	mockRepo.On("UpdateLoanTotalInvested", mock.AnythingOfType("*sql.Tx"), loanID, req1.Amount).Return(updatedLoan1, nil).Once()

	// Mock expectations for second investment
	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(initialLoan, nil).Once()
	mockRepo.On("CreateInvestment", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("*models.Investment")).Return(investment2, nil).Once()
	mockRepo.On("UpdateLoanTotalInvested", mock.AnythingOfType("*sql.Tx"), loanID, req2.Amount).Return(updatedLoan2, nil).Once()

	// Mock transition to invested state (only the second investment will trigger this)
	// The first investment (3000) won't make the loan fully invested (5000 total needed)
	// The second investment (2000) will make it fully invested (3000 + 2000 = 5000)
	mockRepo.On("UpdateLoanState", mock.AnythingOfType("*sql.Tx"), loanID, models.LoanStateInvested).Return(investedLoan, nil).Maybe()
	mockRepo.On("UpdateLoanAgreementLetterURL", mock.AnythingOfType("*sql.Tx"), loanID, mock.AnythingOfType("string")).Return(nil).Maybe()
	mockRepo.On("RecordLoanStateHistory", mock.AnythingOfType("*sql.Tx"), models.LoanStateApproved, investedLoan, mock.AnythingOfType("uuid.UUID"), "Investment target achieved").Return(&models.LoanStateHistory{}, nil).Maybe()

	// Process both investments concurrently with proper synchronization
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []*models.InvestmentResponse
	var errors []error

	wg.Add(2)
	go func() {
		defer wg.Done()
		result, err := service.ProcessInvestment(loanID, req1)
		mu.Lock()
		results = append(results, result)
		errors = append(errors, err)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		result, err := service.ProcessInvestment(loanID, req2)
		mu.Lock()
		results = append(results, result)
		errors = append(errors, err)
		mu.Unlock()
	}()

	wg.Wait()

	// Both investments should succeed
	assert.Len(t, results, 2)
	assert.Len(t, errors, 2)

	for i, err := range errors {
		assert.NoError(t, err, "Investment %d should succeed", i+1)
		assert.NotNil(t, results[i], "Investment %d should return result", i+1)
	}

	// Verify that both investments were processed
	assert.Equal(t, loanID, results[0].LoanID)
	assert.Equal(t, loanID, results[1].LoanID)
	assert.Equal(t, req1.Amount, results[0].Amount)
	assert.Equal(t, req2.Amount, results[1].Amount)

	mockRepo.AssertExpectations(t)
}

func TestLoanService_ProcessDisbursement_Success(t *testing.T) {
	service, mockRepo, mockPayment, _ := setupTestLoanService()

	loanID := uuid.New()
	req := &models.CreateDisbursementRequest{
		FieldOfficerID:          uuid.New(),
		DisbursementDate:        time.Now(),
		SignedAgreementURL:      "https://example.com/agreement.pdf",
		SignedAgreementFileType: models.FileTypePDF,
		DisbursedAmount:         10000.0,
		Notes:                   "Disbursed to borrower",
	}

	loan := createTestLoan(loanID, models.LoanStateInvested, 10000.0)
	loan.AgreementLetterURL = "https://example.com/agreement.pdf"
	disbursement := &models.Disbursement{
		BaseModel:               models.BaseModel{ID: uuid.New()},
		LoanID:                  loanID,
		FieldOfficerID:          req.FieldOfficerID,
		DisbursementDate:        req.DisbursementDate,
		SignedAgreementURL:      req.SignedAgreementURL,
		SignedAgreementFileType: req.SignedAgreementFileType,
		DisbursedAmount:         req.DisbursedAmount,
		Notes:                   req.Notes,
	}
	updatedLoan := createTestLoan(loanID, models.LoanStateDisbursed, 10000.0)
	history := &models.LoanStateHistory{
		BaseModel: models.BaseModel{ID: uuid.New()},
		LoanID:    loanID,
	}
	paymentResult := &adapters.PaymentResult{
		TransactionID: "txn_123",
		Status:        "success",
	}

	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(loan, nil)
	mockRepo.On("CreateDisbursement", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("*models.Disbursement")).Return(disbursement, nil)
	mockRepo.On("UpdateLoanState", mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("uuid.UUID"), models.LoanStateDisbursed).Return(updatedLoan, nil)
	mockRepo.On("RecordLoanStateHistory", mock.AnythingOfType("*sql.Tx"), models.LoanStateInvested, mock.AnythingOfType("*models.Loan"), mock.AnythingOfType("uuid.UUID"), "Loan disbursed").Return(history, nil)
	mockPayment.On("ProcessPayment", 10000.0, mock.AnythingOfType("string")).Return(paymentResult, nil)

	result, err := service.ProcessDisbursement(loanID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, loanID, result.LoanID)
	assert.NotZero(t, result.DisbursedAmount)
	assert.NotEqual(t, uuid.Nil, result.FieldOfficerID)

	mockRepo.AssertExpectations(t)
	mockPayment.AssertExpectations(t)
}

func TestLoanService_ProcessDisbursement_Error(t *testing.T) {
	service, mockRepo, _, _ := setupTestLoanService()

	loanID := uuid.New()
	req := &models.CreateDisbursementRequest{
		FieldOfficerID:          uuid.New(),
		DisbursementDate:        time.Now(),
		SignedAgreementURL:      "https://example.com/agreement.pdf",
		SignedAgreementFileType: models.FileTypePDF,
		DisbursedAmount:         10000.0,
		Notes:                   "Disbursed to borrower",
	}

	// Loan not in invested state - should fail
	loan := createTestLoan(loanID, models.LoanStateApproved, 0)

	mockRepo.On("GetLoanByID", mock.AnythingOfType("*sql.Tx"), loanID).Return(loan, nil)

	result, err := service.ProcessDisbursement(loanID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot transition")

	mockRepo.AssertExpectations(t)
}
