// internal/handlers/loan_handlers.go
package handlers

import (
	"time"

	"loan-service/internal/models"
	"loan-service/internal/services"
	"loan-service/pkg/logger"
	"loan-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoanHandler struct {
	loanService services.LoanServiceInterface
	logger      *logger.Logger
}

func NewLoanHandler(loanService services.LoanServiceInterface, logger *logger.Logger) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
		logger:      logger,
	}
}

// CreateLoan handles loan creation
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req models.CreateLoanRequest

	// First, bind JSON to get the raw data
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate the request using struct tags
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.ValidationErrorFromValidator(c, "Validation failed", err)
		return
	}

	// Business rule validation: ROI cannot be higher than interest rate
	if req.ROI > req.InterestRate {
		response.BadRequest(c, "ROI cannot be higher than interest rate")
		return
	}

	h.logger.Info("Creating loan", map[string]interface{}{
		"request": req,
	})

	loan, err := h.loanService.ProcessCreateLoan(&req)
	if err != nil {
		h.logger.Error("Failed to create loan", map[string]interface{}{
			"error": err.Error(),
		})
		response.BadRequest(c, "Failed to create loan")
		return
	}

	response.Created(c, "Loan created successfully", loan)
}

// GetLoanByID handles getting a loan by ID
func (h *LoanHandler) GetLoanByID(c *gin.Context) {

	loanID := c.Param("loan_id")

	// Parse loan ID
	id, err := uuid.Parse(loanID)
	if err != nil {
		response.BadRequest(c, "Invalid loan ID format")
		return
	}

	loan, err := h.loanService.GetLoanByID(id)
	if err != nil {
		response.BadRequest(c, "Failed to get loan")
		return
	}

	response.Success(c, "Loan retrieved successfully", loan)

}

// ApproveLoan handles loan approval
func (h *LoanHandler) ApproveLoan(c *gin.Context) {

	loanID := c.Param("loan_id")

	// Parse loan ID
	id, err := uuid.Parse(loanID)
	if err != nil {
		response.BadRequest(c, "Invalid loan ID format")
		return
	}

	var req models.CreateApprovalRequest
	req.LoanID = id

	// First, bind JSON to get the raw data
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate the request using struct tags
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.ValidationErrorFromValidator(c, "Validation failed", err)
		return
	}

	// Business rule validation: Approval date should not be in the future
	if req.ApprovalDate.After(time.Now()) {
		response.BadRequest(c, "Approval date cannot be in the future")
		return
	}

	approval, err := h.loanService.ProcessApproveLoan(id, &req)
	if err != nil {
		response.BadRequest(c, "Failed to approve loan")
		return
	}

	response.Success(c, "Loan approved successfully", gin.H{
		"id":                     approval.ID,
		"loan_id":                approval.LoanID,
		"validator_id":           approval.ValidatorID,
		"approval_date":          approval.ApprovalDate,
		"visit_proof_image_url":  approval.VisitProofImageURL,
		"visit_proof_image_type": approval.VisitProofImageType,
		"notes":                  approval.Notes,
		"state":                  models.LoanStateApproved,
	})

}

// AddInvestment handles adding investment to a loan
func (h *LoanHandler) AddInvestment(c *gin.Context) {

	loanID := c.Param("loan_id")

	// Parse loan ID
	id, err := uuid.Parse(loanID)
	if err != nil {
		response.BadRequest(c, "Invalid loan ID format")
		return
	}

	var req models.CreateInvestmentRequest
	req.LoanID = id

	// First, bind JSON to get the raw data
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate the request using struct tags
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.ValidationErrorFromValidator(c, "Validation failed", err)
		return
	}

	// Business rule validation: Investment date should not be in the future
	if req.InvestmentDate.After(time.Now()) {
		response.BadRequest(c, "Investment date cannot be in the future")
		return
	}

	h.logger.Info("Processing investment", map[string]interface{}{
		"loan_id": id.String(),
		"request": req,
	})

	investment, err := h.loanService.ProcessInvestment(id, &req)
	if err != nil {
		h.logger.Error("Failed to process investment", map[string]interface{}{
			"error":   err.Error(),
			"loan_id": id.String(),
		})
		response.BadRequest(c, "Failed to process investment: "+err.Error())
		return
	}

	response.Accepted(c, "Investment processed successfully", investment)
}

// DisburseLoan handles loan disbursement
func (h *LoanHandler) DisburseLoan(c *gin.Context) {

	loanID := c.Param("loan_id")

	// Parse loan ID
	id, err := uuid.Parse(loanID)
	if err != nil {
		response.BadRequest(c, "Invalid loan ID format")
		return
	}

	var req models.CreateDisbursementRequest
	req.LoanID = id

	// First, bind JSON to get the raw data
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate the request using struct tags
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.ValidationErrorFromValidator(c, "Validation failed", err)
		return
	}

	// Business rule validation: Disbursement date should not be in the future
	if req.DisbursementDate.After(time.Now()) {
		response.BadRequest(c, "Disbursement date cannot be in the future")
		return
	}

	h.logger.Info("Processing disbursement", map[string]interface{}{
		"loan_id": id.String(),
		"request": req,
	})

	disbursement, err := h.loanService.ProcessDisbursement(id, &req)
	if err != nil {
		h.logger.Error("Failed to process disbursement", map[string]interface{}{
			"error":   err.Error(),
			"loan_id": id.String(),
		})
		response.BadRequest(c, "Failed to process disbursement: "+err.Error())
		return
	}

	response.Success(c, "Loan disbursed successfully", disbursement)

}
