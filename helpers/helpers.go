package helpers

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetIcon(urlFmt string, iconName string) (image.Image, string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, "", err
	}

	if err := os.MkdirAll(filepath.Join(cacheDir, "lemon"), 0755); err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(cacheDir, "lemon", iconName+".png")
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		if urlFmt != "" {
			if err := DownloadIcon(urlFmt, iconName, filePath); err != nil {
				return nil, "", err
			}
		} else {
			return nil, "", err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, filePath, nil
}

func DownloadIcon(urlFmt string, iconName string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(fmt.Sprintf(urlFmt, iconName))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
