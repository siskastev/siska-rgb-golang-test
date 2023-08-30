package image

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImageToCloudinary(ctx context.Context, imageFile *multipart.FileHeader) (string, error) {

	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", errors.New("failed to initialize cloudinary")
	}

	imageData, err := imageFile.Open()
	if err != nil {
		return "", errors.New("failed to open image file")
	}
	defer imageData.Close()

	result, err := cld.Upload.Upload(ctx, imageData, uploader.UploadParams{
		Folder: os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	})

	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
