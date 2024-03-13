package utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func CheckFileFormat(content []byte) bool {
	fileFormat := http.DetectContentType(content)
	return strings.HasPrefix(fileFormat, "image/png") || strings.HasPrefix(fileFormat, "image/jpeg")
}

func WriteAvatarOnDisk(path string, avatar multipart.File) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = avatar.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, avatar)
	if err != nil {
		return err
	}

	return nil
}
