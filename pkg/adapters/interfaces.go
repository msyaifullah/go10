// pkg/adapters/interfaces.go
package adapters

import (
	"loan-service/internal/models"
	"mime/multipart"
)

type EmailAdapterInterface interface {
	SendEmail(to, subject, body string) error
	GenerateAgreementEmailBody(
		investment *models.Investment,
		loan *models.Loan,
		investor *models.Investor,
		borrower *models.Borrower,
	) string
}

type PaymentAdapterInterface interface {
	ProcessPayment(amount float64, token string) (*PaymentResult, error)
}

type PaymentResult struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type FileAdapterInterface interface {
	UploadFile(file *multipart.FileHeader, entityType string) (*models.FileUpload, error)
}
