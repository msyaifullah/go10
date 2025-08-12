// pkg/logger/gin_middleware.go
package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// maskSensitiveData masks sensitive field values in JSON strings
func maskSensitiveData(data string) string {
	if data == "" {
		return data
	}

	// Common sensitive field patterns - capture field name and mask only the value
	sensitivePatterns := []string{
		`("password"\s*:\s*)"[^"]*"`,
		`("token"\s*:\s*)"[^"]*"`,
		`("api_key"\s*:\s*)"[^"]*"`,
		`("secret"\s*:\s*)"[^"]*"`,
		`("authorization"\s*:\s*)"[^"]*"`,
		`("bearer"\s*:\s*)"[^"]*"`,
		`("credit_card"\s*:\s*)"[^"]*"`,
		`("card_number"\s*:\s*)"[^"]*"`,
		`("ssn"\s*:\s*)"[^"]*"`,
		`("social_security"\s*:\s*)"[^"]*"`,
		`("email"\s*:\s*)"[^"]*"`,
		`("phone"\s*:\s*)"[^"]*"`,
		`("address"\s*:\s*)"[^"]*"`,
		`("account_number"\s*:\s*)"[^"]*"`,
		`("account_name"\s*:\s*)"[^"]*"`,
	}

	maskedData := data
	for _, pattern := range sensitivePatterns {
		re := regexp.MustCompile(pattern)
		maskedData = re.ReplaceAllString(maskedData, `${1}"[**REDACTED by SERVICE**]"`)
	}

	return maskedData
}

// isSensitiveHeader checks if a header name is considered sensitive
func isSensitiveHeader(headerName string) bool {
	sensitiveHeaders := map[string]bool{
		"authorization":   true,
		"x-api-key":       true,
		"x-auth-token":    true,
		"cookie":          true,
		"set-cookie":      true,
		"x-csrf-token":    true,
		"x-forwarded-for": true,
		"x-real-ip":       true,
	}
	// Convert to lowercase for case-insensitive comparison
	return sensitiveHeaders[strings.ToLower(headerName)]
}

type responseWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	headers map[string]string
}

func (r responseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r responseWriter) WriteHeader(statusCode int) {
	// Capture headers before writing status
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	for key, values := range r.Header() {
		if len(values) > 0 {
			r.headers[key] = values[0]
		}
	}
	r.ResponseWriter.WriteHeader(statusCode)
}

func GinMiddleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Read and restore request body for logging
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Wrap response writer to capture response body
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
			headers:        make(map[string]string),
		}
		c.Writer = w

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)

		logData := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       path,
			"status":     c.Writer.Status(),
			"duration":   duration.String(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// Add request headers (will be masked if sensitive)
		if len(c.Request.Header) > 0 {
			headers := make(map[string]string)
			for key, values := range c.Request.Header {
				if len(values) > 0 {
					// Mask sensitive header values directly
					headerValue := values[0] // Take first value if multiple
					if isSensitiveHeader(key) {
						headerValue = "[**REDACTED by SERVICE**]"
					}
					headers[key] = headerValue
				}
			}
			logData["request_headers"] = headers
		}

		// Add request_id if present in context
		if requestID, exists := c.Get("request_id"); exists {
			logData["request_id"] = requestID
		}

		// Add request body if present (will be masked if sensitive)
		if len(bodyBytes) > 0 && len(bodyBytes) < 1024 { // Only log small bodies
			maskedBody := maskSensitiveData(string(bodyBytes))
			// Try to parse as JSON to avoid double escaping
			var jsonData interface{}
			if err := json.Unmarshal([]byte(maskedBody), &jsonData); err == nil {
				logData["request_body"] = jsonData
			} else {
				logData["request_body"] = maskedBody
			}
		}

		// Add response body if it's not too large
		if w.body.Len() > 0 && w.body.Len() < 1024 {
			maskedBody := maskSensitiveData(w.body.String())
			// Try to parse as JSON to avoid double escaping
			var jsonData interface{}
			if err := json.Unmarshal([]byte(maskedBody), &jsonData); err == nil {
				logData["response_body"] = jsonData
			} else {
				logData["response_body"] = maskedBody
			}
		}

		// Add response headers (will be masked if sensitive)
		if len(w.headers) > 0 {
			maskedHeaders := make(map[string]string)
			for key, value := range w.headers {
				// Mask sensitive header values directly
				if isSensitiveHeader(key) {
					maskedHeaders[key] = "[**REDACTED by SERVICE**]"
				} else {
					maskedHeaders[key] = value
				}
			}
			logData["response_headers"] = maskedHeaders
		}

		// Add error if present
		if len(c.Errors) > 0 {
			logData["errors"] = c.Errors.Errors()
		}

		if c.Writer.Status() >= 400 {
			logger.Error("HTTP request completed with error", logData)
		} else {
			logger.Info("HTTP request completed", logData)
		}
	}
}
