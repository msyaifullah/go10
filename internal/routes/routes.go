// internal/routes/routes.go
package routes

import (
	"loan-service/internal/application"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

/**
KSUID (K-Sortable Unique IDentifier) is a great choice for request IDs because:
	It's time-sortable (useful for logging and debugging)
	It's more compact than UUID
	It's collision-resistant
	It's URL-safe
*/

// RequestIDMiddleware generates a unique request ID for tracking
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = ksuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func SetupRoutes(r *gin.Engine, app *application.Application) {
	// Add request ID middleware
	r.Use(RequestIDMiddleware())

	// Health check endpoint
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":      "healthy",
			"service":     app.Config.App.Name,
			"environment": app.Config.App.Environment,
		})
	})

	api := r.Group("/api/v1")

	// Loan routes
	loans := api.Group("/loans")
	{
		// Create loan
		loans.POST("/", app.LoanHandler.CreateLoan)
		loans.GET("/:loan_id", app.LoanHandler.GetLoanByID)
		loans.POST("/:loan_id/approve", app.LoanHandler.ApproveLoan)
		loans.POST("/:loan_id/invest", app.LoanHandler.AddInvestment)
		loans.POST("/:loan_id/disburse", app.LoanHandler.DisburseLoan)
	}

	// File upload routes
	files := api.Group("/files")
	{
		files.POST("/upload", app.FileHandler.UploadFile)
	}
}
