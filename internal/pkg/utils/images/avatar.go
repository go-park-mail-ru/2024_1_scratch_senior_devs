package images

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func CheckFileFormat(content []byte) string {
	fileFormat := http.DetectContentType(content)

	if strings.HasPrefix(fileFormat, "image/png") {
		return ".png"
	}

	if strings.HasPrefix(fileFormat, "image/jpeg") {
		return ".jpeg"
	}

	return ""
}

func WriteAvatarOnDisk(path string, avatar io.ReadSeeker) error {
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
