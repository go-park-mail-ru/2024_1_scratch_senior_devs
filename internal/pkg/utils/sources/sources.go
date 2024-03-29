package sources

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func CheckFormat(choice map[string]string, content []byte) string {
	fileFormat := http.DetectContentType(content)

	for mimeTime, format := range choice {
		if strings.HasPrefix(fileFormat, mimeTime) {
			return format
		}
	}

	return ""
}

func WriteFileOnDisk(path string, resource io.ReadSeeker) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = resource.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resource)
	if err != nil {
		return err
	}

	return nil
}
