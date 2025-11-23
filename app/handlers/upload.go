package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"knowledge-capsule-api/pkg/utils"
)

// UploadHandler handles file uploads.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPost) {
		return
	}

	// Limit upload size to 10MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	// Create uploads directory if not exists
	uploadDir := "data/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Return file URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	utils.JSONResponse(w, http.StatusCreated, true, "File uploaded successfully", map[string]string{
		"file_url": fileURL,
	})
}
