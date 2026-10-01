package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var allowedImageExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".webp": true,
}

// uploadImage accepts a multipart form file under the field name "file",
// saves it into ./uploads with a unique name, and returns its public URL.
// Point any *_url field (logo_url, media_image_url, icon_url, banner_image_url…)
// at the returned url.
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field \"file\" is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type — use png, jpg, jpeg, gif, or webp"})
		return
	}

	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not prepare uploads folder"})
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadsDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + filename})
}
