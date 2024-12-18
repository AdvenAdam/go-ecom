package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var Validate = validator.New()

const MAX_UPLOAD_SIZE = 10 * 1024 * 1024 // 10MB

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")

	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})

}

// UploadFile uploads the given file to the given destination directory.
// It returns the uploaded file path or an error.
//
// The destination directory will be created if it doesn't exist.
// The file will be closed after uploading.
func UploadFile(file multipart.File, filename, dest string) (string, error) {
	// Create destination directory if it doesn't exist
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create destination file
	dst, err := os.Create(filepath.Join(dest, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filepath.Join(dest, filename), nil
}

// GenerateUniqueFilename generates a unique filename with the given extension.
func GenerateUniqueFilename(ext string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String() + ext, nil
}

func GetFileNameDestination(fileHeader *multipart.FileHeader) (string, string, error) {
	// Define the upload destination based on a string parameter
	filename, err := GenerateUniqueFilename(filepath.Ext(fileHeader.Filename))
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	// TODO : MOVE FOLDER DESTINATION TO ENV
	uploadDestination := "uploads/images/" + filename

	return filename, uploadDestination, nil
}
