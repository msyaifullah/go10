// pkg/adapters/email_adapter.go
package adapters

import (
	"fmt"
	"net/smtp"

	"loan-service/internal/models"
	"loan-service/pkg/config"
	"loan-service/pkg/logger"
)

type EmailAdapter struct {
	config config.EmailConfig
	logger *logger.Logger
}

func NewEmailAdapter(cfg config.EmailConfig, logger *logger.Logger) EmailAdapterInterface {
	return &EmailAdapter{
		config: cfg,
		logger: logger,
	}
}

func (a *EmailAdapter) SendEmail(to, subject, body string) error {
	a.logger.Debug("Sending email", map[string]interface{}{
		"to":       to,
		"subject":  subject,
		"provider": a.config.Provider,
	})

	switch a.config.Provider {
	case "smtp":
		return a.sendSMTP(to, subject, body)
	case "sendgrid":
		return a.sendSendGrid(to, subject, body)
	case "mock":
		return a.sendMock(to, subject, body)
	default:
		return fmt.Errorf("unsupported email provider: %s", a.config.Provider)
	}
}

func (a *EmailAdapter) sendSMTP(to, subject, body string) error {
	auth := smtp.PlainAuth("", a.config.SMTPUsername, a.config.SMTPPassword, a.config.SMTPHost)

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	addr := fmt.Sprintf("%s:%d", a.config.SMTPHost, a.config.SMTPPort)

	err := smtp.SendMail(addr, auth, a.config.FromAddress, []string{to}, []byte(msg))
	if err != nil {
		a.logger.Error("Failed to send SMTP email", map[string]interface{}{
			"to":    to,
			"error": err.Error(),
		})
		return err
	}

	a.logger.Info("SMTP email sent successfully", map[string]interface{}{
		"to":      to,
		"subject": subject,
	})

	return nil
}

func (a *EmailAdapter) sendSendGrid(to, subject, body string) error {
	// Implementation for SendGrid API
	a.logger.Info("SendGrid email sent (mock)", map[string]interface{}{
		"to":      to,
		"subject": subject,
		"api_key": a.config.APIKey,
	})
	return nil
}

func (a *EmailAdapter) sendMock(to, subject, body string) error {
	a.logger.Info("Mock email sent", map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	})
	return nil
}

// GenerateAgreementEmailBody generates the email body for investment agreements
func (a *EmailAdapter) GenerateAgreementEmailBody(
	investment *models.Investment,
	loan *models.Loan,
	investor *models.Investor,
	borrower *models.Borrower,
) string {
	return fmt.Sprintf(`Dear %s,

Your investment in Loan #%s has been successfully processed.

Investment Details:
- Investment Amount: Rp %.2f
- Expected Return: Rp %.2f
- Investment Date: %s

Loan Details:
- Borrower: %s
- Principal Amount: Rp %.2f
- Interest Rate: %.2f%%
- ROI: %.2f%%

Agreement Letter: %s

Please review the agreement letter and contact us if you have any questions.

Best regards,
Go10 Team`,
		investor.Name,
		loan.ID.String()[:8],
		investment.Amount,
		investment.ExpectedReturn,
		investment.InvestmentDate.Format("January 2, 2006"),
		borrower.FullName(),
		loan.PrincipalAmount,
		loan.InterestRate*100,
		loan.ROI*100,
		loan.AgreementLetterURL,
	)
}
