package utils

import (
	"bytes"
	"image/jpeg"
	"io"
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
			_, err := buf.ReadFrom(part)
			if err != nil {
				part.Close()
				return err
			}
			data["file"] = buf
			data["fileName"] = part.FileName()
			part.Close()
		} else if part.FormName() == "image" {
			image, err := jpeg.Decode(part)
			if err != nil {
				return err
			}
			data["image"] = image
			part.Close()
		} else {
			for key, _ := range data {
				if part.FormName() == key {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(part)
					if err != nil {
						part.Close()
						return err
					}
					data[key] = buf.String()
					part.Close()
					break
				}
			}
		}
	}

	return nil
}
