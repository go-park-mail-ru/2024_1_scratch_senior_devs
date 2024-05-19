package exportpdf

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/satori/uuid"
	"strings"
	"unicode/utf8"
)

const maxFilenameLength = 50

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
	return fmt.Sprintf(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"/><title>%s</title><style>%s</style></head><body>%s</body></html>`, getNoteTitle(basicHTML), stylesLight, basicHTML)
}

func processImg(document *goquery.Document) map[uuid.UUID]int {
	document.Find("img").Not("[data-imgid]").Remove()

	pictureCount := 0
	result := make(map[uuid.UUID]int, 0)
	document.Find("img[data-imgid]").Each(func(i int, s *goquery.Selection) {
		imgID, exists := s.Attr("data-imgid")
		if exists {
			pictureCount++
			result[uuid.FromStringOrNil(imgID)] = pictureCount
			s.ReplaceWithHtml(fmt.Sprintf(`<div>----- картинка %d -----</div>`, pictureCount))
		}
	})

	return result
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

func processFile(document *goquery.Document) map[uuid.UUID]string {
	filenames := make(map[uuid.UUID]string, 0)
	document.Find("button[data-fileid]").Each(func(i int, s *goquery.Selection) {
		attachName, exists := s.Attr("data-filename")
		if utf8.RuneCountInString(attachName) > maxFilenameLength {
			attachName = attachName[:maxFilenameLength]
		}

		if exists {
			s.Find(".file-name").Each(func(i int, selection *goquery.Selection) {
				selection.SetText(attachName)
			})

			attachID, found := s.Attr("data-fileid")
			if found {
				filenames[uuid.FromStringOrNil(attachID)] = attachName
			}

			format := `<a href="https://you-note.ru">`
			s.Contents().Each(func(i int, contentSelection *goquery.Selection) {
				inner, err := contentSelection.Html()
				if err != nil {
					return
				}
				format += inner
			})
			format += "</a>"

			docA, err := goquery.NewDocumentFromReader(strings.NewReader(format))
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

	return filenames
}

func processIframe(document *goquery.Document) {
	document.Find("iframe").Each(func(i int, s *goquery.Selection) {
		iframeSrc, exists := s.Attr("src")
		if exists {
			a := fmt.Sprintf(`<a href="%s">%s</a>`, iframeSrc, iframeSrc)
			s.BeforeHtml(a)
		}
		s.Remove()
	})
}

func prepareHTML(basicHTML string) (string, map[uuid.UUID]int, map[uuid.UUID]string, error) {
	noteHTML := wrap(basicHTML)

	document, err := goquery.NewDocumentFromReader(strings.NewReader(noteHTML))
	if err != nil {
		return noteHTML, map[uuid.UUID]int{}, map[uuid.UUID]string{}, errors.New("invalid input HTML")
	}

	picturesOrder := processImg(document)
	processSubnote(document)
	filenames := processFile(document)
	processIframe(document)

	newHTML, err := document.Html()
	if err != nil {
		return noteHTML, map[uuid.UUID]int{}, map[uuid.UUID]string{}, errors.New("internal error while parsing processed HTML")
	}

	return newHTML, picturesOrder, filenames, nil
}
