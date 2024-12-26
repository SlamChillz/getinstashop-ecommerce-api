package utils

import (
	"fmt"
	"log"
	"path/filepath"
)

func GenerateMigrationFilePath() string {
	absPath, err := filepath.Abs(".")
	if err != nil {
		log.Println("Error getting the absolute path:", err)
		return ""
	}
	relativePath := "internal/db/migrations"
	fullPath := filepath.Join(absPath, relativePath)
	fileURL := fmt.Sprintf("file://%s", fullPath)
	return fileURL
}
