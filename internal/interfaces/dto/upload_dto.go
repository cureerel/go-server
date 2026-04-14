// internal/interfaces/dto/upload_dto.go
package dto

type UploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type DeleteUploadRequest struct {
	Key string `json:"key" binding:"required"`
}
