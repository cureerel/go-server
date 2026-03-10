// internal/interfaces/http/handler/upload_handler.go
package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	storageinfra "github.com/cureerel/gotemplate/internal/infrastructure/storage"
	"github.com/cureerel/gotemplate/internal/interfaces/dto"
	"github.com/cureerel/gotemplate/pkg/apperror"
	"github.com/gin-gonic/gin"
)

const (
	maxUploadBytes  = 5 * 1024 * 1024
	uploadFormField = "image"
)

var allowedMIMEs = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

type UploadHandler struct {
	storage storageinfra.Provider
}

func NewUploadHandler(storage storageinfra.Provider) *UploadHandler {
	return &UploadHandler{storage: storage}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(maxUploadBytes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large — max 5 MB"})
		return
	}

	file, header, err := c.Request.FormFile(uploadFormField)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'image' missing from form"})
		return
	}
	defer file.Close()

	if header.Size > maxUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large — max 5 MB"})
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxUploadBytes+1))
	if err != nil {
		respondErr(c, apperror.NewInternal(err, "failed to read upload"))
		return
	}
	if int64(len(data)) > maxUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large — max 5 MB"})
		return
	}

	mime := http.DetectContentType(data)
	mime = strings.SplitN(mime, ";", 2)[0]
	if !allowedMIMEs[mime] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "unsupported file type",
			"allowed": "jpeg, png, webp, gif",
		})
		return
	}

	folder := c.DefaultQuery("folder", "general")
	if !validFolder(folder) {
		folder = "general"
	}

	uid, _ := getUID(c)
	key := fmt.Sprintf("%d-%s", uid, sanitiseFilename(header.Filename))

	result, err := h.storage.Upload(c.Request.Context(), storageinfra.UploadInput{
		Key:         key,
		Body:        strings.NewReader(string(data)),
		ContentType: mime,
		Folder:      folder,
	})
	if err != nil {
		respondErr(c, apperror.NewInternal(err, "upload failed"))
		return
	}

	respondCreated(c, dto.UploadResponse{URL: result.URL, Key: result.Key})
}

func (h *UploadHandler) DeleteImage(c *gin.Context) {
	var req dto.DeleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, _ := getUID(c)
	if !hasRole(c, "admin") {
		expectedPrefix := fmt.Sprintf("%d-", uid)
		base := filepath.Base(req.Key)
		if !strings.HasPrefix(base, expectedPrefix) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't own this file"})
			return
		}
	}

	if err := h.storage.Delete(c.Request.Context(), req.Key); err != nil {
		respondErr(c, apperror.NewInternal(err, "delete failed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

func validFolder(f string) bool {
	switch f {
	case "blogs", "services", "avatars", "general":
		return true
	}
	return false
}

func sanitiseFilename(name string) string {
	name = filepath.Base(name)
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}