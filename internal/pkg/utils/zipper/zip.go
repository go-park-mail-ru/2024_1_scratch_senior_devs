package zipper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/satori/uuid"
	"io"
	"os"
	"path"
	"strings"
)

// CreateYouNoteZip создает zip-архив, в котором находится PDF-файл заметки и файлы вложений
func CreateYouNoteZip(pdfTitle string, pdfBytes []byte, paths []string, username string, picturesOrder map[uuid.UUID]int, filenames map[uuid.UUID]string) (*bytes.Buffer, error) {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	defer zipWriter.Close()

	if err := addFileToZip(zipWriter, pdfTitle, pdfBytes); err != nil {
		return zipBuffer, err
	}

	for _, filepath := range paths {
		if err := addFileToZipFromPath(zipWriter, filepath, picturesOrder, filenames); err != nil {
			return zipBuffer, err
		}
	}

	if err := addFileToZip(zipWriter, "YouNote❤️.txt", []byte(fmt.Sprintf("Hello, %s!\nHope you like it :D", username))); err != nil {
		return zipBuffer, err
	}

	return zipBuffer, nil
}

func addFileToZip(zipWriter *zip.Writer, filename string, content []byte) error {
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = fileWriter.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func addFileToZipFromPath(zipWriter *zip.Writer, filename string, picturesOrder map[uuid.UUID]int, filenames map[uuid.UUID]string) error {
	file, err := os.Open(path.Join(os.Getenv("ATTACHES_BASE_PATH"), filename))
	if err != nil {
		return nil // skip this attach if it wasn`t found
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	fileExt := path.Ext(filename)
	fileID := uuid.FromStringOrNil(strings.TrimSuffix(filename, fileExt))

	place, found := picturesOrder[fileID]
	if found {
		header.Name = fmt.Sprintf("картинка %d%s", place, fileExt)
	}

	name, found := filenames[fileID]
	if found {
		header.Name = name
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}
