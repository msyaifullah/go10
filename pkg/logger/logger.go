// pkg/logger/logger.go
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"loan-service/pkg/config"
)

// LoggerInterface defines the contract for loggers
type LoggerInterface interface {
	Info(message string, data map[string]interface{})
	Error(message string, data map[string]interface{})
	Debug(message string, data map[string]interface{})
	Warn(message string, data map[string]interface{})
	Fatal(message string, data map[string]interface{})
}

type Logger struct {
	config        config.LoggerConfig
	maskSensitive bool
	sensitiveKeys []string
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func NewLogger(cfg config.LoggerConfig) *Logger {
	return &Logger{
		config:        cfg,
		maskSensitive: cfg.MaskSensitive,
		sensitiveKeys: []string{
			"password", "token", "secret", "key", "auth", "credential",
			"api_key", "secret_key", "webhook_secret", "smtp_password",
			"authorization", "bearer", "session", "cookie",
		},
	}
}

func (l *Logger) maskData(data map[string]interface{}) map[string]interface{} {
	if !l.maskSensitive {
		return data
	}

	masked := make(map[string]interface{})
	for k, v := range data {
		masked[k] = l.maskValue(k, v)
	}
	return masked
}

func (l *Logger) maskValue(key string, value interface{}) interface{} {
	keyLower := strings.ToLower(key)

	// Check if key contains sensitive terms
	for _, sensitive := range l.sensitiveKeys {
		if strings.Contains(keyLower, sensitive) {
			return l.maskString(fmt.Sprintf("%v", value))
		}
	}

	// Check if value is a string that looks like sensitive data
	if str, ok := value.(string); ok {
		// Mask potential tokens, API keys, etc.
		if l.isSensitiveString(str) {
			return l.maskString(str)
		}
	}

	// Check if value is a map (nested data) and recursively mask it
	if nestedMap, ok := value.(map[string]interface{}); ok {
		return l.maskData(nestedMap)
	}

	// Check if value is a slice and mask each item
	if slice, ok := value.([]interface{}); ok {
		maskedSlice := make([]interface{}, len(slice))
		for i, item := range slice {
			if itemMap, ok := item.(map[string]interface{}); ok {
				maskedSlice[i] = l.maskData(itemMap)
			} else {
				maskedSlice[i] = l.maskValue("", item)
			}
		}
		return maskedSlice
	}

	return value
}

func (l *Logger) isSensitiveString(s string) bool {
	// Pattern for API keys, tokens, etc.
	patterns := []string{
		`^sk_[a-zA-Z0-9_]+$`,         // Stripe secret keys
		`^pk_[a-zA-Z0-9_]+$`,         // Stripe public keys
		`^SG\.[a-zA-Z0-9_\-]+$`,      // SendGrid API keys
		`^[A-Za-z0-9+/]{32,}={0,2}$`, // Base64 encoded strings (potential tokens)
		`^[a-zA-Z0-9]{32,}$`,         // Long alphanumeric strings (potential keys)
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, s); matched {
			return true
		}
	}

	return false
}

func (l *Logger) maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

func (l *Logger) log(level, message string, data map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Data:      l.maskData(data),
	}

	if l.config.Format == "json" {
		jsonData, _ := json.Marshal(entry)
		log.Println(string(jsonData))
	} else {
		dataStr := ""
		if len(entry.Data) > 0 {
			// Use MarshalIndent for better formatting without escaped quotes
			dataBytes, _ := json.MarshalIndent(entry.Data, "", "  ")
			dataStr = fmt.Sprintf("\nData: %s", string(dataBytes))
		}
		log.Printf("[%s] %s: %s%s", entry.Level, entry.Timestamp, entry.Message, dataStr)
	}
}

func (l *Logger) Debug(message string, data map[string]interface{}) {
	if l.shouldLog("debug") {
		l.log("DEBUG", message, data)
	}
}

func (l *Logger) Info(message string, data map[string]interface{}) {
	if l.shouldLog("info") {
		l.log("INFO", message, data)
	}
}

func (l *Logger) Warn(message string, data map[string]interface{}) {
	if l.shouldLog("warn") {
		l.log("WARN", message, data)
	}
}

func (l *Logger) Error(message string, data map[string]interface{}) {
	if l.shouldLog("error") {
		l.log("ERROR", message, data)
	}
}

func (l *Logger) Fatal(message string, data map[string]interface{}) {
	l.log("FATAL", message, data)
	os.Exit(1)
}

func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
		"fatal": 4,
	}

	configLevel, exists := levels[l.config.Level]
	if !exists {
		configLevel = 1 // default to info
	}

	logLevel, exists := levels[level]
	if !exists {
		return false
	}

	return logLevel >= configLevel
}
