package net

import (
	"io"
	"net/http"
	"os"
)

func GetHtml(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if bytes, err := io.ReadAll(res.Body); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

// Download downloads file from given url to filename
func Download(url string, filename string) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = output.Close()
	}()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	_, err = io.Copy(output, res.Body)
	return err
}
