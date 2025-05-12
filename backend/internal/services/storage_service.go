package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageService defines methods for file storage operations
type StorageService interface {
	UploadFile(tenantID, fileType string, file *multipart.FileHeader) (string, error)
	GetFileURL(tenantID, fileKey string) (string, error)
	DeleteFile(tenantID, fileKey string) error
}

// LocalStorageService implements StorageService for local file system
type LocalStorageService struct {
	basePath string
}

// NewStorageService creates a new storage service based on configuration
func NewStorageService(config map[string]string) (StorageService, error) {
	// Default to local storage
	basePath := config["basePath"]
	if basePath == "" {
		basePath = "./data/uploads" // Default base path
	}

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &LocalStorageService{
		basePath: basePath,
	}, nil
}

// UploadFile uploads a file to local storage
func (s *LocalStorageService) UploadFile(tenantID, fileType string, file *multipart.FileHeader) (string, error) {
	// Create tenant directory if it doesn't exist
	tenantDir := filepath.Join(s.basePath, tenantID, fileType)
	if err := os.MkdirAll(tenantDir, 0755); err != nil {
		return "", err
	}

	// Generate a unique file path
	filename := filepath.Base(file.Filename)
	timestamp := time.Now().Unix()
	fileKey := fmt.Sprintf("%s/%s/%d_%s", tenantID, fileType, timestamp, filename)
	filePath := filepath.Join(s.basePath, fileKey)

	// Open the source file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the file contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return fileKey, nil
}

// GetFileURL gets the file path for local storage
func (s *LocalStorageService) GetFileURL(tenantID, fileKey string) (string, error) {
	// Ensure the file key belongs to the correct tenant
	if !strings.HasPrefix(fileKey, tenantID+"/") {
		return "", fmt.Errorf("unauthorized access to file")
	}

	filePath := filepath.Join(s.basePath, fileKey)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}

	// For local storage, return the absolute file path
	// In a production environment with a web server, this would be a URL
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// DeleteFile deletes a file from local storage
func (s *LocalStorageService) DeleteFile(tenantID, fileKey string) error {
	// Ensure the file key belongs to the correct tenant
	if !strings.HasPrefix(fileKey, tenantID+"/") {
		return fmt.Errorf("unauthorized access to file")
	}

	filePath := filepath.Join(s.basePath, fileKey)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
