package cldnary

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ClientCloudinary struct {
	Client *cloudinary.Cloudinary
	Folder string
}

func NewCloudinary(apiUrl string, folder string) (*ClientCloudinary, error) {
	client, err := cloudinary.NewFromURL(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cloudinary, error: %v", err)
	}
	return &ClientCloudinary{
		Client: client,
		Folder: folder,
	}, nil
}

func (c *ClientCloudinary) UploadImage(ctx context.Context, image *multipart.FileHeader, folder string) (string, string, error) {
	src, err := image.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	folder = fmt.Sprintf("%s-%s", c.Folder, folder)
	result, err := c.Client.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: folder,
	})

	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}

func (c *ClientCloudinary) DeleteImage(ctx context.Context, publicID string) error {
	_, err := c.Client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return err
	}

	return nil
}
