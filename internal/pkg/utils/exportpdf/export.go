package exportpdf

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/satori/uuid"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GeneratePDF(basicHTML string) ([]byte, string, map[uuid.UUID]int, map[uuid.UUID]string, error) {
	noteHTML, picturesOrder, filenames, err := prepareHTML(basicHTML)
	if err != nil {
		return nil, "", picturesOrder, filenames, err
	}

	t := time.Now().Unix()

	if _, err := os.Stat("cloneHTML/"); os.IsNotExist(err) {
		if err := os.Mkdir("cloneHTML/", 0777); err != nil {
			return nil, "", picturesOrder, filenames, errors.New("1: " + err.Error())
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, "", picturesOrder, filenames, errors.New("2: " + err.Error())
	}
	defer os.RemoveAll(dir + "/cloneHTML")

	if err := os.WriteFile("cloneHTML/"+strconv.FormatInt(t, 10)+".html", []byte(noteHTML), 0644); err != nil { //nolint
		return nil, "", picturesOrder, filenames, errors.New("3: " + err.Error())
	}

	file, err := os.Open("cloneHTML/" + strconv.FormatInt(t, 10) + ".html")
	if err != nil {
		return nil, "", picturesOrder, filenames, errors.New("4: " + err.Error())
	}
	defer file.Close()

	wkhtmltopdf.SetPath("/bin/wkhtmltopdf")
	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, "", picturesOrder, filenames, errors.New("5: " + err.Error())
	}

	generator.NoPdfCompression.Set(true)
	generator.AddPage(wkhtmltopdf.NewPageReader(file))
	generator.PageSize.Set(wkhtmltopdf.PageSizeA4)
	generator.Dpi.Set(300)

	if err = generator.Create(); err != nil {
		return nil, "", picturesOrder, filenames, errors.New("6: " + err.Error())
	}

	return generator.Bytes(), getNoteTitle(basicHTML), picturesOrder, filenames, nil
}
