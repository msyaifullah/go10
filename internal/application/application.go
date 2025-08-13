// internal/application/application.go
package application

import (
	"database/sql"
	"errors"

	"loan-service/internal/handlers"
	"loan-service/internal/repositories"
	"loan-service/internal/services"
	"loan-service/pkg/adapters"
	"loan-service/pkg/config"
	"loan-service/pkg/logger"
	"loan-service/pkg/redis"
)

type Application struct {
	// Core dependencies
	Config *config.Config
	Logger *logger.Logger
	DB     *sql.DB
	Redis  *redis.RedisClient

	// Repositories
	LoanRepo repositories.LoanRepositoryInterface

	// Adapters
	EmailAdapter   adapters.EmailAdapterInterface
	PaymentAdapter adapters.PaymentAdapterInterface
	FileAdapter    adapters.FileAdapterInterface

	// Services
	LoanService services.LoanServiceInterface
	CronService *services.CronService

	// Handlers
	LoanHandler *handlers.LoanHandler
	FileHandler *handlers.FileHandler
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) WithConfig(cfg *config.Config) *Application {
	app.Config = cfg
	return app
}

func (app *Application) WithLogger(logger *logger.Logger) *Application {
	app.Logger = logger
	return app
}

func (app *Application) WithDatabase(db *sql.DB) *Application {
	app.DB = db
	return app
}

func (app *Application) WithRedis() *Application {
	redisClient, err := redis.NewConnection(app.Config.Redis, app.Logger)
	if err != nil {
		app.Logger.Error("Failed to initialize Redis", map[string]interface{}{
			"error": err.Error(),
		})
		// You might want to handle this error differently based on your requirements
		// For now, we'll just log the error and continue
	}
	app.Redis = redisClient
	return app
}

func (app *Application) WithRepositories() *Application {
	app.LoanRepo = repositories.NewLoanRepository(app.DB, app.Logger)
	return app
}

func (app *Application) WithAdapters() *Application {
	app.EmailAdapter = adapters.NewEmailAdapter(app.Config.Email, app.Logger)
	app.PaymentAdapter = adapters.NewPaymentAdapter(app.Config.Payment, app.Logger)
	app.FileAdapter = adapters.NewFileAdapter(app.Logger)
	return app
}

func (app *Application) WithServices() *Application {
	app.LoanService = services.NewLoanService(
		app.LoanRepo,
		app.PaymentAdapter,
		app.EmailAdapter,
		app.Logger,
		app.DB,
	)

	app.CronService = services.NewCronService(
		app.LoanRepo,
		app.EmailAdapter,
		app.Logger,
		app.DB,
		app.Config,
	)

	return app
}

func (app *Application) WithHandlers() *Application {
	app.LoanHandler = handlers.NewLoanHandler(app.LoanService, app.Logger)
	app.FileHandler = handlers.NewFileHandler(app.Logger, app.FileAdapter)
	return app
}

func (app *Application) Validate() error {
	if app.Config == nil {
		return errors.New("config not initialized")
	}
	if app.Logger == nil {
		return errors.New("logger not initialized")
	}
	if app.DB == nil {
		return errors.New("database not initialized")
	}
	if app.Redis == nil {
		return errors.New("redis not initialized")
	}

	return nil
}
