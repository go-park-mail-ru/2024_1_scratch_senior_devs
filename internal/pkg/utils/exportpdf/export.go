package exportpdf

import (
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GeneratePDF(basicHTML string) ([]byte, string, error) {
	noteHTML, err := prepareHTML(basicHTML)
	if err != nil {
		return nil, "", err
	}

	t := time.Now().Unix()

	if _, err := os.Stat("cloneHTML/"); os.IsNotExist(err) {
		if err := os.Mkdir("cloneHTML/", 0777); err != nil {
			return nil, "", err
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(dir + "/cloneHTML")

	if err := os.WriteFile("cloneHTML/"+strconv.FormatInt(t, 10)+".html", []byte(noteHTML), 0644); err != nil {
		return nil, "", err
	}

	file, err := os.Open("cloneHTML/" + strconv.FormatInt(t, 10) + ".html")
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, "", err
	}

	generator.NoPdfCompression.Set(true)
	generator.AddPage(wkhtmltopdf.NewPageReader(file))
	generator.PageSize.Set(wkhtmltopdf.PageSizeA4)
	generator.Dpi.Set(300)

	if err = generator.Create(); err != nil {
		return nil, "", err
	}

	return generator.Bytes(), getNoteTitle(basicHTML), nil
}
