package yes24

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/librarios/cli/internal/pkg/net"
	"os"
	"regexp"
	"strings"
)

func getGoodsHtml(id string) (string, error) {
	url := fmt.Sprintf("https://www.yes24.com/Product/Goods/%s", id)
	return net.GetHtml(url)
}

type GoodsInfo struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Genre       string        `json:"genre"`
	Keywords    string        `json:"keywords"`
	Author      NameWithType  `json:"author"`
	Publisher   NameWithType  `json:"publisher"`
	Url         string        `json:"url"`
	WorkExample []WorkExample `json:"workExample"`
}

type NameWithType struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

type WorkExample struct {
	Isbn            string          `json:"isbn"`
	DatePublished   string          `json:"datePublished"`
	PotentialAction PotentialAction `json:"potentialAction"`
}

type PotentialAction struct {
	ExpectsAcceptanceOf ExpectsAcceptanceOf `json:"expectsAcceptanceOf"`
}

type ExpectsAcceptanceOf struct {
	Type          string `json:"@type"`
	Price         string `json:"price"`
	PriceCurrency string `json:"priceCurrency"`
	Availability  string `json:"availability"`
}

func extractGoodsInfo(html string) (*GoodsInfo, error) {
	re := regexp.MustCompile(`(?m)<script type="application/ld\+json">((\n|.)*?)</script>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) == 0 {
		return nil, errors.New("goods info not foundin html")
	}

	// parse goods info JSON
	match := matches[1]
	var info GoodsInfo
	if err := json.Unmarshal([]byte(match), &info); err != nil {
		return nil, err
	}

	return &info, nil
}

type BookInfo struct {
	Title         string
	Isbn          string
	Author        string
	PublishedDate string
	Publisher     string
	Price         string
	PriceCurrency string
}

func (s *BookInfo) GetPublishedYear() string {
	return s.PublishedDate[:4]
}

func GetBookInfo(id string) (*BookInfo, error) {
	html, err := getGoodsHtml(id)
	if err != nil {
		return nil, err
	}

	info, err := extractGoodsInfo(html)
	if err != nil {
		return nil, err
	}

	bookInfo := BookInfo{
		Title:     info.Name,
		Author:    info.Author.Name,
		Publisher: info.Publisher.Name,
	}
	if len(info.WorkExample) > 0 {
		work := info.WorkExample[0]
		bookInfo.Isbn = work.Isbn
		bookInfo.PublishedDate = work.DatePublished
		bookInfo.Price = work.PotentialAction.ExpectsAcceptanceOf.Price
		bookInfo.PriceCurrency = work.PotentialAction.ExpectsAcceptanceOf.PriceCurrency
	}

	return &bookInfo, nil
}

func WriteBookInfo(info *BookInfo, filename string) error {
	lines := []string{
		fmt.Sprintf("isbn=%s", info.Isbn),
		"origPubDate=",
		"origTitle=",
		fmt.Sprintf("pubDate=%s", info.PublishedDate),
		"scanDate=",
		"scanPages=",
		fmt.Sprintf("price=%s", info.Price),
	}
	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0777)
}
