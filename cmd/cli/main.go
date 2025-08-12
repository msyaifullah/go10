package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	port           int
	host           string
	debug          bool
	migrationsPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go10svc",
	Short: "go10 CLI for project's utilities",
}

// migrateCmd represents the database migration command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database migrations...")
		if debug {
			fmt.Println("Debug mode enabled")
		}

		// Build database URL
		dbHost := viper.GetString("database.host")
		dbPort := viper.GetInt("database.port")
		dbName := viper.GetString("database.name")
		dbUser := viper.GetString("database.user")
		dbPassword := viper.GetString("database.password")

		// Use DATABASE_URL from environment if available, otherwise build from config
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				dbUser, dbPassword, dbHost, dbPort, dbName)
		}

		// Get migrations path
		migrationsPath := viper.GetString("migration.path")
		if migrationsPath == "" {
			migrationsPath = "migrations"
		}

		// Ensure migrations directory exists
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			fmt.Printf("Creating migrations directory: %s\n", migrationsPath)
			if err := os.MkdirAll(migrationsPath, 0755); err != nil {
				fmt.Printf("Error creating migrations directory: %v\n", err)
				return
			}
		}

		// Create source URL for file-based migrations
		sourceURL := fmt.Sprintf("file://%s", filepath.Join(migrationsPath))

		if debug {
			fmt.Printf("Database URL: %s\n", databaseURL)
			fmt.Printf("Source URL: %s\n", sourceURL)
		}

		// Initialize migrate instance
		m, err := migrate.New(sourceURL, databaseURL)
		if err != nil {
			fmt.Printf("Error initializing migrations: %v\n", err)
			return
		}
		defer m.Close()

		// Run migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Error running migrations: %v\n", err)
			return
		}

		if err == migrate.ErrNoChange {
			fmt.Println("No pending migrations to apply.")
		} else {
			fmt.Println("Migrations completed successfully!")
		}

		// Get current version
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			fmt.Printf("Error getting migration version: %v\n", err)
		} else if err != migrate.ErrNilVersion {
			fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		}
	},
}

// migrateDownCmd represents the database migration down command
var migrateDownCmd = &cobra.Command{
	Use:   "migrate-down",
	Short: "Rollback database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rolling back database migrations...")
		if debug {
			fmt.Println("Debug mode enabled")
		}

		// Build database URL
		dbHost := viper.GetString("database.host")
		dbPort := viper.GetInt("database.port")
		dbName := viper.GetString("database.name")
		dbUser := viper.GetString("database.user")
		dbPassword := viper.GetString("database.password")

		// Use DATABASE_URL from environment if available, otherwise build from config
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				dbUser, dbPassword, dbHost, dbPort, dbName)
		}

		// Get migrations path
		migrationsPath := viper.GetString("migration.path")
		if migrationsPath == "" {
			migrationsPath = "migrations"
		}

		// Create source URL for file-based migrations
		sourceURL := fmt.Sprintf("file://%s", filepath.Join(migrationsPath))

		if debug {
			fmt.Printf("Database URL: %s\n", databaseURL)
			fmt.Printf("Source URL: %s\n", sourceURL)
		}

		// Initialize migrate instance
		m, err := migrate.New(sourceURL, databaseURL)
		if err != nil {
			fmt.Printf("Error initializing migrations: %v\n", err)
			return
		}
		defer m.Close()

		// Get current version before rolling back
		currentVersion, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			fmt.Printf("Error getting current migration version: %v\n", err)
			return
		}
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations have been applied yet.")
			return
		}
		if dirty {
			fmt.Println("Warning: Database is in dirty state. Consider fixing manually.")
		}

		fmt.Printf("Current migration version: %d\n", currentVersion)

		// Rollback one migration
		if err := m.Steps(-1); err != nil {
			fmt.Printf("Error rolling back migration: %v\n", err)
			return
		}

		// Get new version after rollback
		newVersion, _, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			fmt.Printf("Error getting new migration version: %v\n", err)
		} else if err == migrate.ErrNilVersion {
			fmt.Println("All migrations have been rolled back.")
		} else {
			fmt.Printf("Successfully rolled back to version: %d\n", newVersion)
		}
	},
}

// envCmd represents the environment preparation command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Prepare environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Preparing environment variables...")
		if debug {
			fmt.Println("Debug mode enabled")
		}

		// Check required environment variables
		requiredVars := []string{"DATABASE_URL", "API_KEY", "ENVIRONMENT"}
		missingVars := []string{}

		for _, varName := range requiredVars {
			if os.Getenv(varName) == "" {
				missingVars = append(missingVars, varName)
			}
		}

		if len(missingVars) > 0 {
			fmt.Printf("Missing required environment variables: %v\n", missingVars)
			fmt.Println("Please set the following environment variables:")
			for _, varName := range missingVars {
				fmt.Printf("  export %s=<value>\n", varName)
			}
		} else {
			fmt.Println("All required environment variables are set!")
		}

		fmt.Printf("Configuration loaded from: %s\n", viper.ConfigFileUsed())
	},
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Printf("  Database Host: %s\n", viper.GetString("database.host"))
		fmt.Printf("  Database Port: %d\n", viper.GetInt("database.port"))
		fmt.Printf("  Database Name: %s\n", viper.GetString("database.name"))
		fmt.Printf("  Debug: %t\n", viper.GetBool("debug"))
		fmt.Printf("  Config File: %s\n", viper.ConfigFileUsed())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in configs directory
		viper.AddConfigPath("configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Set default values
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "app_db")
	viper.SetDefault("database.user", "app_user")
	viper.SetDefault("database.password", "app_password")
	viper.SetDefault("migration.path", "migrations")
	viper.SetDefault("debug", false)

	// Read environment variables
	viper.AutomaticEnv() // read in environment variables that match
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is configs/cli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")

	// Local flags for migrate command
	migrateCmd.Flags().StringVarP(&host, "host", "H", "localhost", "database host")
	migrateCmd.Flags().IntVarP(&port, "port", "p", 5432, "database port")
	migrateCmd.Flags().StringVarP(&migrationsPath, "path", "m", "migrations", "migrations directory path")

	// Bind flags to viper
	viper.BindPFlag("database.host", migrateCmd.Flags().Lookup("host"))
	viper.BindPFlag("database.port", migrateCmd.Flags().Lookup("port"))
	viper.BindPFlag("migration.path", migrateCmd.Flags().Lookup("path"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add commands
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateDownCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(configCmd)
}

func main() {
	Execute()
}
