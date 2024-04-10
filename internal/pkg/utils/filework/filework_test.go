package filework

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"testing"
)

func TestGetFormat(t *testing.T) {
	choice := map[string]string{
		"audio/mp4":     ".mp4",
		"audio/mpeg":    ".mp3",
		"audio/vnd.wav": ".wav",
		"image/gif":     ".gif",
		"image/jpeg":    ".jpeg",
		"image/png":     ".png",
		"video/mp4":     ".mp4",
	}

	tests := []struct {
		name    string
		content []byte
		result  string
	}{
		{
			name:    "Test_CheckFormat_Jpeg",
			content: []byte{255, 216, 255, 224},
			result:  ".jpeg",
		},
		{
			name:    "Test_CheckFormat_Unknown",
			content: []byte{0, 1, 2, 3},
			result:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFormat(choice, tt.content)

			assert.Equal(t, tt.result, result)
		})
	}
}

func TestSaveImageAsWebp(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	blue := color.RGBA{0, 0, 255, 255}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, blue)
		}
	}

	err := SaveImageAsWebp("test_image", img, quality)
	assert.Nil(t, err)

	_, err = os.Stat("test_image.webp")
	assert.False(t, os.IsNotExist(err))

	err = os.Remove("test_image.webp")
	assert.Nil(t, err)
}

func TestSaveFile(t *testing.T) {
	t.Run("Successful file saving", func(t *testing.T) {
		testData := []byte("test data")
		resource := bytes.NewReader(testData)

		err := SaveFile("testfile", ".txt", resource)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		fileData, err := os.ReadFile("testfile.txt")
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}

		if !bytes.Equal(fileData, testData) {
			t.Errorf("File content mismatch, expected %s, got %s", testData, fileData)
		}

		if err := os.Remove("testfile.txt"); err != nil {
			fmt.Println(err.Error())
		}
	})

	t.Run("Error seeking resource", func(t *testing.T) {
		resource := &mockReadSeeker{readErr: errors.New("seek error")}
		err := SaveFile("testfile", ".txt", resource)
		if err == nil {
			t.Errorf("Expected error seeking resource")
		}

		if err := os.Remove("testfile.txt"); err != nil {
			fmt.Println(err.Error())
		}
	})

	t.Run("Error copying resource to file", func(t *testing.T) {
		resource := &mockReadSeeker{readErr: errors.New("copy error")}
		err := SaveFile("testfile", ".txt", resource)
		if err == nil {
			t.Errorf("Expected error copying resource to file")
		}

		if err := os.Remove("testfile.txt"); err != nil {
			fmt.Println(err.Error())
		}
	})
}

type mockReadSeeker struct {
	readErr error
}

func (m *mockReadSeeker) Read(p []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	return len(p), nil
}

func (m *mockReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	return 0, nil
}

func TestWriteFileOnDisk(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		oldExtension string
		resource     io.ReadSeeker
		expectedExt  string
		isErr        bool
	}{
		{
			name:         "Test_WriteFileOnDisk_Jpg",
			path:         "test.jpg",
			oldExtension: ".jpg",
			resource:     getTestImage("jpg"),
			expectedExt:  ".webp",
			isErr:        false,
		},
		{
			name:         "Test_WriteFileOnDisk_Png",
			path:         "test.png",
			oldExtension: ".png",
			resource:     getTestImage("png"),
			expectedExt:  ".webp",
			isErr:        false,
		},
		{
			name:        "Test_WriteFileOnDisk_Other",
			path:        "test.txt",
			resource:    getTestImage("invalid"),
			expectedExt: "",
			isErr:       false,
		},
	}

	for _, tt := range tests {
		actualExt, err := WriteFileOnDisk(tt.path, tt.oldExtension, tt.resource)

		assert.Equal(t, tt.expectedExt, actualExt)

		if tt.isErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		if err := os.Remove(tt.path + tt.expectedExt); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func getTestImage(format string) io.ReadSeeker {
	var b bytes.Buffer

	switch format {
	case "jpg":
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		_ = jpeg.Encode(&b, img, nil)
	case "png":
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		_ = png.Encode(&b, img)
	default:
		b.WriteString("not_an_image")
	}

	return bytes.NewReader(b.Bytes())
}
