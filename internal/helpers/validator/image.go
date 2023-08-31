package validator

import (
	"mime/multipart"
	"strings"
)

func IsValidImageType(file *multipart.FileHeader) bool {
	validTypes := []string{"image/jpeg", "image/jpg", "image/png"}

	contentType := file.Header.Get("Content-Type")
	for _, validType := range validTypes {
		if strings.EqualFold(contentType, validType) {
			return true
		}
	}

	return false
}

func IsValidSizeImage(file *multipart.FileHeader) bool {
	size := file.Size
	if size > 5*1024*1024 {
		return false
	}
	return true
}
