package utils

import (
	"fmt"
	"io"
	"io/ioutil"
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

func LoadLocalFIle(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
