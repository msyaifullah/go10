// pkg/adapters/email_adapter.go
package adapters

import (
	"fmt"
	"loan-service/internal/models"
	"loan-service/pkg/logger"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileAdapter struct {
	logger *logger.Logger
}

func NewFileAdapter(logger *logger.Logger) FileAdapterInterface {
	return &FileAdapter{
		logger: logger,
	}
}

func (a *FileAdapter) UploadFile(file *multipart.FileHeader, entityType string) (*models.FileUpload, error) {
	a.logger.Debug("Uploading file", map[string]interface{}{
		"file_name":    file.Filename,
		"file_size":    file.Size,
		"content_type": file.Header.Get("Content-Type"),
		"entity_type":  entityType,
	})

	// Extract file extension and determine file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	var fileExtensionType models.FileType
	switch ext {
	case ".pdf":
		fileExtensionType = models.FileTypePDF
	case ".jpg", ".jpeg":
		fileExtensionType = models.FileTypeJPEG
	case ".png":
		fileExtensionType = models.FileTypePNG
	default:
		fileExtensionType = models.FileTypePDF // Default fallback
	}

	// Generate filename based on entity type
	var filename string
	uuidStr := uuid.New().String()[:8]

	switch entityType {
	case "disbursement":
		filename = fmt.Sprintf("agreements/signed_%s%s", uuidStr, ext)
	case "approval":
		filename = fmt.Sprintf("proof/visit_%s%s", uuidStr, ext)
	case "loan":
		filename = fmt.Sprintf("agreements/loan_%s%s", uuidStr, ext)
	default:
		filename = fmt.Sprintf("file_%s%s", uuidStr, ext)
	}

	// Generate file path and URL
	filePath := "/uploads/" + filename
	fileURL := "https://storage.example.com" + filePath

	// Create FileUpload model
	fileUpload := &models.FileUpload{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FileName:    filename,
		FileType:    fileExtensionType,
		FileSize:    file.Size,
		FilePath:    filePath,
		FileURL:     fileURL,
		ContentType: file.Header.Get("Content-Type"),
		UploadedBy:  uuid.Nil,   // This should be set by the handler based on authenticated user
		EntityType:  entityType, // Use the entityType parameter from form
		EntityID:    uuid.Nil,   // This should be set by the handler based on context
		IsActive:    true,
	}

	return fileUpload, nil
}
