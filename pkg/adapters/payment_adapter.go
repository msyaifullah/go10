// pkg/adapters/payment_adapter.go
package adapters

import (
	"fmt"

	"loan-service/pkg/config"
	"loan-service/pkg/logger"
)

type PaymentAdapter struct {
	config config.PaymentConfig
	logger *logger.Logger
}

func NewPaymentAdapter(cfg config.PaymentConfig, logger *logger.Logger) PaymentAdapterInterface {
	return &PaymentAdapter{
		config: cfg,
		logger: logger,
	}
}

func (a *PaymentAdapter) ProcessPayment(amount float64, token string) (*PaymentResult, error) {
	a.logger.Debug("Processing payment", map[string]interface{}{
		"amount":   amount,
		"token":    token,
		"provider": a.config.Provider,
	})

	switch a.config.Provider {
	case "stripe":
		return a.processStripe(amount, token)
	case "mock":
		return a.processMock(amount, token)
	default:
		return nil, fmt.Errorf("unsupported payment provider: %s", a.config.Provider)
	}
}

func (a *PaymentAdapter) processStripe(amount float64, token string) (*PaymentResult, error) {
	// Implementation for Stripe API
	a.logger.Info("Stripe payment processed (mock)", map[string]interface{}{
		"amount":     amount,
		"token":      token,
		"secret_key": a.config.SecretKey,
	})

	return &PaymentResult{
		TransactionID: "stripe_txn_123456",
		Status:        "success",
		Message:       "Payment processed successfully via Stripe",
	}, nil
}

func (a *PaymentAdapter) processMock(amount float64, token string) (*PaymentResult, error) {
	a.logger.Info("Mock payment processed", map[string]interface{}{
		"amount": amount,
		"token":  token,
	})

	return &PaymentResult{
		TransactionID: "mock_txn_123456",
		Status:        "success",
		Message:       "Payment processed successfully (mock)",
	}, nil
}
