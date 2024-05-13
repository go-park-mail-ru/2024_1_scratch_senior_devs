package exportpdf

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"path/filepath"
	"strings"
)

func getNoteTitle(basicHTML string) string {
	title := ""

	document, err := goquery.NewDocumentFromReader(strings.NewReader(basicHTML))
	if err != nil {
		return title
	}

	document.Find(".note-title").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})

	return title
}

func wrap(basicHTML string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"/><title>%s</title><style>%s</style></head><body>%s</body></html>`, getNoteTitle(basicHTML), styles, basicHTML)
}

func processImg(document *goquery.Document) {
	document.Find("img[data-imgid]").Each(func(i int, s *goquery.Selection) {
		imgID, exists := s.Attr("data-imgid")
		if exists {
			newSrc := fmt.Sprintf("https://you-note.ru/attaches/%s.webp", imgID)
			s.SetAttr("src", newSrc)
		}
	})

	document.Find("img").Not("[data-imgid]").Remove()
}

func processSubnote(document *goquery.Document) {
	document.Find("button[data-noteid]").Each(func(i int, s *goquery.Selection) {
		noteID, exists := s.Attr("data-noteid")
		if exists {
			format := `<a href="https://you-note.ru/notes/%s">`
			s.Contents().Each(func(i int, contentSelection *goquery.Selection) {
				inner, err := contentSelection.Html()
				if err != nil {
					return
				}
				format += inner
			})
			format += "</a>"

			docA, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf(format, noteID)))
			if err != nil {
				return
			}
			a := docA.Find("a")

			a.Each(func(i int, contentSelection *goquery.Selection) {
				className, exists := s.Attr("class")
				if exists {
					a.SetAttr("class", className)
				}

				contenteditable, exists := s.Attr("contenteditable")
				if exists {
					a.SetAttr("contenteditable", contenteditable)
				}
			})

			s.ReplaceWithNodes(a.Get(0))
		}
	})
}

func processFile(document *goquery.Document) {
	document.Find("button[data-fileid]").Each(func(i int, s *goquery.Selection) {
		attachID, exists := s.Attr("data-fileid")
		if exists {
			format := `<a href="https://you-note.ru/attaches/%s%s">`
			s.Contents().Each(func(i int, contentSelection *goquery.Selection) {
				inner, err := contentSelection.Html()
				if err != nil {
					return
				}
				format += inner
			})
			format += "</a>"

			extension := ""
			dataFilename, exists := s.Attr("data-filename")
			if exists {
				extension = filepath.Ext(dataFilename)
			}

			docA, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf(format, attachID, extension)))
			if err != nil {
				return
			}
			a := docA.Find("a")

			a.Each(func(i int, contentSelection *goquery.Selection) {
				className, exists := s.Attr("class")
				if exists {
					a.SetAttr("class", className)
				}

				contenteditable, exists := s.Attr("contenteditable")
				if exists {
					a.SetAttr("contenteditable", contenteditable)
				}

				dataFilename, exists := s.Attr("data-filename")
				if exists {
					a.SetAttr("data-filename", dataFilename)
				}
			})

			s.ReplaceWithNodes(a.Get(0))
		}
	})
}

func processIframe(document *goquery.Document) {
	document.Find("iframe").Each(func(i int, s *goquery.Selection) {
		iframeSrc, exists := s.Attr("src")
		if exists {
			a := fmt.Sprintf(`<a href="%s">%s</a>`, iframeSrc, iframeSrc)
			s.BeforeHtml(a)
		}
	})
}

func prepareHTML(basicHTML string) (string, error) {
	noteHTML := wrap(basicHTML)

	document, err := goquery.NewDocumentFromReader(strings.NewReader(noteHTML))
	if err != nil {
		return noteHTML, errors.New("invalid input HTML")
	}

	processImg(document)
	processSubnote(document)
	processFile(document)
	processIframe(document)

	newHTML, err := document.Html()
	if err != nil {
		return noteHTML, errors.New("internal error while parsing processed HTML")
	}

	return newHTML, nil
}
