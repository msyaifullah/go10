package models

import (
	"time"

	"github.com/google/uuid"
)

// EmailNotification tracks email notifications sent to investors
type EmailNotification struct {
	BaseModel
	InvestorID   uuid.UUID  `json:"investor_id" validate:"required"`
	LoanID       uuid.UUID  `json:"loan_id" validate:"required"`
	EmailType    string     `json:"email_type" validate:"required"` // agreement_notification, etc.
	EmailSubject string     `json:"email_subject" validate:"required"`
	EmailBody    string     `json:"email_body" validate:"required"`
	SentAt       time.Time  `json:"sent_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	OpenedAt     *time.Time `json:"opened_at,omitempty"`
	Status       string     `json:"status" validate:"required"` // sent, delivered, opened, failed
	ErrorMessage string     `json:"error_message,omitempty"`

	// Relationships
	Investor *Investor `json:"investor,omitempty"`
	Loan     *Loan     `json:"loan,omitempty"`
}

// FileUpload tracks uploaded files for audit and management
type FileUpload struct {
	BaseModel
	FileName    string    `json:"file_name" validate:"required"`
	FileType    FileType  `json:"file_type" validate:"required"`
	FileSize    int64     `json:"file_size" validate:"required"`
	FilePath    string    `json:"file_path" validate:"required"`
	FileURL     string    `json:"file_url" validate:"required"`
	ContentType string    `json:"content_type" validate:"required"`
	UploadedBy  uuid.UUID `json:"uploaded_by" validate:"required"`
	EntityType  string    `json:"entity_type" validate:"required"` // loan, approval, disbursement
	EntityID    uuid.UUID `json:"entity_id" validate:"required"`
	IsActive    bool      `json:"is_active"`
}
