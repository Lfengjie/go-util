package util

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(writer *multipart.Writer, fieldname, filename, mime string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	if mime == "" {
		mime = "application/octet-stream"
	}
	h.Set("Content-Type", mime)
	return writer.CreatePart(h)
}

func makeFormData(filename, mime string, content io.Reader) (formData io.Reader, contentType string, err error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := createFormFile(writer, "filecontent", filename, mime)
	if err != nil {
		return
	}
	_, err = io.Copy(part, content)
	if err != nil {
		return
	}

	formData = buf
	contentType = writer.FormDataContentType()
	writer.Close()

	return
}

func PutObjectToGiftv2(url, filename string, s string) error {
	reader := strings.NewReader(s)
	mime := "application/json"

	formData, mime, err := makeFormData(filename, mime, reader)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, formData)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", mime)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			err = fmt.Errorf("status is %d", resp.StatusCode)
			return err
		}
	}
	return nil
}
