// pkg/response/response.go
package response

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response represents the unified API response structure
type Response struct {
	Status  string      `json:"status"`           // "success" or "failed"
	Message string      `json:"message"`          // message to client
	Code    string      `json:"code"`             // error code or success code
	Data    interface{} `json:"data"`             // object or array of objects
	Errors  []string    `json:"errors,omitempty"` // validation errors list
}

// PaginationResponse represents paginated data response
type PaginationResponse struct {
	Response
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`        // current page number
	Limit      int   `json:"limit"`       // items per page
	TotalItems int64 `json:"total_items"` // total number of items
	TotalPages int   `json:"total_pages"` // total number of pages
	HasNext    bool  `json:"has_next"`    // has next page
	HasPrev    bool  `json:"has_prev"`    // has previous page
}

// Success codes
const (
	CodeSuccess   = "SUCCESS"
	CodeCreated   = "CREATED"
	CodeUpdated   = "UPDATED"
	CodeDeleted   = "DELETED"
	CodeProcessed = "PROCESSED"
	CodeAccepted  = "ACCEPTED"
)

// Error codes
const (
	CodeBadRequest      = "BAD_REQUEST"
	CodeNotFound        = "NOT_FOUND"
	CodeInternalError   = "INTERNAL_ERROR"
	CodeValidationError = "VALIDATION_ERROR"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
)

// Success responses
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Code:    CodeSuccess,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: message,
		Code:    CodeCreated,
		Data:    data,
	})
}

func Updated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Code:    CodeUpdated,
		Data:    data,
	})
}

func Deleted(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Code:    CodeDeleted,
	})
}

func Processed(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Code:    CodeProcessed,
	})
}

func Accepted(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusAccepted, Response{
		Status:  "success",
		Message: message,
		Code:    CodeAccepted,
		Data:    data,
	})
}

// Error responses
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Status:  "failed",
		Message: message,
		Code:    CodeBadRequest,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Status:  "failed",
		Message: message,
		Code:    CodeNotFound,
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Status:  "failed",
		Message: message,
		Code:    CodeInternalError,
	})
}

func ValidationError(c *gin.Context, message string, errors ...string) {
	response := Response{
		Status:  "failed",
		Message: message,
		Code:    CodeValidationError,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	c.JSON(http.StatusBadRequest, response)
}

func ValidationErrorFromValidator(c *gin.Context, message string, err error) {
	validationErrors := ParseValidationErrors(err)
	ValidationError(c, message, validationErrors...)
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Status:  "failed",
		Message: message,
		Code:    CodeUnauthorized,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Status:  "failed",
		Message: message,
		Code:    CodeForbidden,
	})
}

// Pagination helpers
func PaginatedSuccess(c *gin.Context, message string, data interface{}, page, limit int, totalItems int64) {
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := PaginationInfo{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	c.JSON(http.StatusOK, PaginationResponse{
		Response: Response{
			Status:  "success",
			Message: message,
			Code:    CodeSuccess,
			Data:    data,
		},
		Pagination: pagination,
	})
}

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) (page, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ = strconv.Atoi(pageStr)
	limit, _ = strconv.Atoi(limitStr)

	// Ensure minimum values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit to prevent abuse
	}

	return page, limit
}

// ParseValidationErrors converts validator.ValidationErrors to a list of readable error messages
func ParseValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			// Get the field name (remove struct name prefix)
			fieldName := fieldError.Field()
			if strings.Contains(fieldName, ".") {
				parts := strings.Split(fieldName, ".")
				fieldName = parts[len(parts)-1] // Get the last part after the dot
			}

			// Convert field name to snake_case for better readability
			fieldName = strings.ToLower(fieldName)

			// Create user-friendly error messages
			switch fieldError.Tag() {
			case "required":
				errors = append(errors, fieldName+" is required")
			case "gt":
				errors = append(errors, fieldName+" must be greater than "+fieldError.Param())
			case "gte":
				errors = append(errors, fieldName+" must be greater than or equal to "+fieldError.Param())
			case "lt":
				errors = append(errors, fieldName+" must be less than "+fieldError.Param())
			case "lte":
				errors = append(errors, fieldName+" must be less than or equal to "+fieldError.Param())
			case "email":
				errors = append(errors, fieldName+" must be a valid email address")
			case "uuid":
				errors = append(errors, fieldName+" must be a valid UUID")
			default:
				errors = append(errors, fieldName+" failed validation: "+fieldError.Tag())
			}
		}
	} else {
		// If it's not a validation error, just return the error message
		errors = append(errors, err.Error())
	}

	return errors
}
