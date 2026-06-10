package utils

import (
	"bytes"
	"io"
	"strings"

	"github.com/dslipak/pdf"
)

func ExtractPDFText(file io.Reader, size int64) (string, error) {
	buf, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	reader, err := pdf.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
