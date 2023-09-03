package yes24

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/librarios/cli/internal/pkg/net"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getPreviewHtml(id string) (string, error) {
	url := fmt.Sprintf("https://www.yes24.com/Product/Viewer/Preview/%s", id)
	return net.GetHtml(url)
}

type PreviewInput struct {
	Pages struct {
		Page []PagePreview `json:"PAGE"`
	}
	PageDomain string `json:"pagedomain"`
	GoodsNo    string `json:"goodsno"`
}

type PagePreview struct {
	GoodsNo      int32
	OrderNo      int8
	OriginalName string
	SmallImage   PreviewImage
	MiddleImage  PreviewImage
	LargeImage   PreviewImage
	Bookmark     string
}

type PreviewImage struct {
	Name     string
	Width    int16
	Height   int16
	FileSize int16
}

func extractPreviewImageUrls(html string) ([]string, error) {
	// find preview definition block
	re := regexp.MustCompile(`(?m)var input = ({(\n|.)*?});`)
	matches := re.FindStringSubmatch(html)
	if len(matches) == 0 {
		return nil, errors.New("page preview info not found in html")
	}

	// parse preview definition JSON
	inputJson := matches[1]
	inputJson = strings.ReplaceAll(inputJson, "'isonepagemode': onepagemode,", "")
	inputJson = strings.ReplaceAll(inputJson, "'", "\"")

	var input PreviewInput
	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return nil, err
	}

	// create preview image URL list
	var urls []string
	for _, preview := range input.Pages.Page {
		url := fmt.Sprintf("%s/%s", input.PageDomain, preview.LargeImage.Name)
		urls = append(urls, url)
	}
	return urls, nil
}

func DownloadPreviewImages(id string, dir string) error {
	html, err := getPreviewHtml(id)
	if err != nil {
		return err
	}

	urls, err := extractPreviewImageUrls(html)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	for i, url := range urls {
		ext := filepath.Ext(url)
		outputFilename := filepath.Join(dir, fmt.Sprintf("%03d%s", i, ext))
		if err = net.Download(url, outputFilename); err != nil {
			return fmt.Errorf("failed to download %s : %+v", url, err)
		} else {
			slog.Info("downloaded", "filename", outputFilename)
		}
	}

	return nil
}
