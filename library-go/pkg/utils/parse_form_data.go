package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

// ParseMultiPartFormData - for parsing send map with key of multipart and value type (text or file)
// if file included also added a new key to map ["fileName"] that contains the file name
func ParseMultiPartFormData(r *http.Request, data map[string]interface{}) error {

	part, err := r.MultipartReader()
	if err != nil {
		return err
	}
	for {
		part, err := part.NextPart()
		if err == io.EOF {
			break
		}
		if part.FormName() == "file" {
			buf := new(bytes.Buffer)
			data["file"] = buf
			data["fileName"] = part.FileName()
		} else {
			for key, _ := range data {
				if part.FormName() == key {
					buf := new(bytes.Buffer)
					buf.ReadFrom(part)
					data[key] = buf.String()
					break
				}
			}
		}
	}

	return nil
}

func ParseMultiPartFormData2(part *multipart.Reader, data map[string]interface{}) error {

	for {
		part, err := part.NextPart()
		if err == io.EOF {
			break
		}
		if part.FormName() == "file" {
			buf := new(bytes.Buffer)
			bytesRead, err := buf.ReadFrom(part)
			if err != nil {
				return err
			}
			if bytesRead < 1 {
				return io.EOF
			}
			data["file"] = buf
			data["fileName"] = part.FileName()
		} else {
			for key, _ := range data {
				if part.FormName() == key {
					buf := new(bytes.Buffer)
					buf.ReadFrom(part)
					data[key] = buf.String()
					break
				}
			}
		}
	}

	return nil
}
