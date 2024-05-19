package zipper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// CreateYouNoteZip создает zip-архив, в котором находится PDF-файл заметки и файлы вложений
func CreateYouNoteZip(zipFileName string, pdfTitle string, pdfBytes []byte, paths []string, username string) (*bytes.Buffer, error) {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	defer zipWriter.Close()

	if err := addFileToZip(zipWriter, pdfTitle, pdfBytes); err != nil {
		return zipBuffer, err
	}

	for _, path := range paths {
		if err := addFileToZipFromPath(zipWriter, path); err != nil {
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

func addFileToZipFromPath(zipWriter *zip.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
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
