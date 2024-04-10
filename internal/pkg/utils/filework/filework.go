package filework

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

const quality = 80

func GetFormat(choice map[string]string, content []byte) string {
	fileFormat := http.DetectContentType(content)

	for mimeType, format := range choice {
		if strings.HasPrefix(fileFormat, mimeType) {
			return format
		}
	}

	return ""
}

// SaveImageAsWebp godoc
// compresses img and writes to disk as <path.webp>
func SaveImageAsWebp(path string, img image.Image, quality float32) error {
	file, err := os.Create(path + ".webp")
	if err != nil {
		return err
	}
	defer file.Close()

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return err
	}

	if err := webp.Encode(file, img, options); err != nil {
		return err
	}

	return nil
}

// SaveFile godoc
// writes resource to disk as <path+extension>
func SaveFile(path string, extension string, resource io.ReadSeeker) error {
	file, err := os.Create(path + extension)
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

// WriteFileOnDisk godoc
// Converts jpg/jpeg/png to webp and saves to disk
// If cant convert to webp - just saves to disk
func WriteFileOnDisk(path string, oldExtension string, resource io.ReadSeeker) (string, error) {
	_, err := resource.Seek(0, 0)
	if err != nil {
		return "", err
	}

	var img image.Image
	img, _, err = image.Decode(resource)
	if err != nil {
		if err = SaveFile(path, oldExtension, resource); err != nil {
			return "", err
		}

		return oldExtension, nil
	}

	if err = SaveImageAsWebp(path, img, quality); err != nil {
		return "", err
	}

	return ".webp", nil
}
