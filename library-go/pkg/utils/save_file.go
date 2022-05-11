package utils

import (
	"fmt"
	"io"
	"os"
)

func SaveFile(path string, fileName string, fileSrc io.Reader) error {

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	fileLocalPath := fmt.Sprintf("%s/%s", path, fileName)
	file, err := os.Create(fileLocalPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, fileSrc)
	if err != nil {
		return err
	}

	return nil
}
