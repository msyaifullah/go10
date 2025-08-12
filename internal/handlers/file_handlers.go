package handlers

import (
	"loan-service/pkg/adapters"
	"loan-service/pkg/logger"
	"loan-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	logger      *logger.Logger
	fileAdapter adapters.FileAdapterInterface
}

func NewFileHandler(logger *logger.Logger, fileAdapter adapters.FileAdapterInterface) *FileHandler {
	return &FileHandler{
		logger:      logger,
		fileAdapter: fileAdapter,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("Failed to get multipart form", map[string]interface{}{
			"error": err.Error(),
		})
		response.BadRequest(c, "Failed to get multipart form")
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		h.logger.Error("No files uploaded", map[string]interface{}{})
		response.BadRequest(c, "No files uploaded")
		return
	}

	// Extract entity type from form
	entityType := "file" // Default value
	if entityTypes := form.Value["entity_type"]; len(entityTypes) > 0 {
		entityType = entityTypes[0]
	}

	// Validate entity type
	validEntityTypes := map[string]bool{
		"approval":     true,
		"disbursement": true,
		"loan":         true,
	}
	if !validEntityTypes[entityType] {
		h.logger.Error("Invalid entity type", map[string]interface{}{
			"entity_type": entityType,
		})
		response.BadRequest(c, "Invalid entity type. Must be 'approval', 'disbursement', or 'loan'")
		return
	}

	file := files[0]

	fileUpload, err := h.fileAdapter.UploadFile(file, entityType)
	if err != nil {
		h.logger.Error("Failed to upload file", map[string]interface{}{
			"error": err.Error(),
		})
		response.BadRequest(c, "Failed to upload file")
		return
	}

	response.Success(c, "File uploaded successfully", fileUpload)
}
