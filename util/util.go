package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
)

var (
	replacer *strings.Replacer = strings.NewReplacer(
		":", "_",
		"<", "[",
		">", "]",
		"|", "-",
		"\"", "",
		"/", ".",
		"\\", ".",
		"?", "",
		"*", "",
	)
)

func CleanString(str string) string {
	var b strings.Builder
	for _, r := range str {
		if r >= 32 && r <= 126 {
			b.WriteRune(r)
		}
	}

	return replacer.Replace(b.String())
}

func ReadImageFromResp(resp *http.Response) ([]byte, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errStr := fmt.Sprintf("failed to read body: %s", err.Error())
		return nil, errors.New(errStr)
	}

	return data, nil
}

// * Note: panel might be moved to another folder
type Panel struct {
	URL       *url.URL
	Image     []byte
	ImageType string
	Width     float64
	Height    float64
}

func CreatePDF(folder string, title string, pages []Panel) error {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "in",
		Size:           gofpdf.SizeType{Wd: 8.33, Ht: 13.33}, // desired comic size is 800x1280 pixels which convert to "inches" is 8.33x13.33
	})

	pdf.SetMargins(0.0, 0.0, 0.0)
	pdf.SetCellMargin(0.0)

	for i, p := range pages {
		pdf.AddPageFormat("P", gofpdf.SizeType{Wd: p.Width, Ht: p.Height})
		pdf.RegisterImageReader(title+strconv.Itoa(i), p.ImageType, bytes.NewBuffer(p.Image))
		if pdf.Ok() {
			options := gofpdf.ImageOptions{
				ReadDpi:   false,
				ImageType: p.ImageType,
			}

			pdf.ImageOptions(title+strconv.Itoa(i), 0, pdf.GetY(), p.Width, p.Height, false, options, 0, "")
		}
	}

	title = CleanString(title)
	err := pdf.OutputFileAndClose("./" + folder + "/" + title + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func AddLog(msg string, folder string) {
	err := ioutil.WriteFile("./"+folder+"/log.txt", []byte(msg), 0655)
	if err != nil {
		// try log again?
	}
}
