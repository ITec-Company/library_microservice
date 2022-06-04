package utils

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
)

type Format string
type Extension string

var (
	FormatOriginal Format = "original"
	FormatQVGA     Format = "QVGA"
	FormatVGA      Format = "VGA"
	FormatHD720p   Format = "HD720p"

	JPG Extension = ".jpg"

	ErrNilPointer = errors.New("nil pointer reference")
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

func LoadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func SaveImage(image *image.Image, path string) error {

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	imagesMap, err := ResizeImage(image)
	if err != nil {
		return err
	}

	for format, img := range imagesMap {
		file, err := os.Create(fmt.Sprintf("%s/%s.jpg", path, string(format)))
		if err != nil {
			return err
		}

		if err = jpeg.Encode(file, img, nil); err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func ResizeImage(original *image.Image) (map[Format]image.Image, error) {
	if original == nil {
		return nil, ErrNilPointer
	}
	images := make(map[Format]image.Image)
	images[FormatOriginal] = *original
	images[FormatQVGA] = resize.Resize(320, 240, *original, resize.Lanczos3)
	images[FormatVGA] = resize.Resize(640, 480, *original, resize.Lanczos3)
	images[FormatHD720p] = resize.Resize(1280, 720, *original, resize.Lanczos3)

	return images, nil
}

func GetImageFromLocalStore(path string, format Format, extension Extension) (*image.Image, error) {
	file, err := os.Open(fmt.Sprintf("%s%s%s", path, string(format), string(extension)))
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
