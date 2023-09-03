package yes24

import (
	"fmt"
	"os"
	"strings"
)

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
