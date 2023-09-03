package yes24

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/librarios/cli/internal/pkg/net"
	"html"
	"regexp"
)

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

func getGoodsHtml(id string) (string, error) {
	url := fmt.Sprintf("https://www.yes24.com/Product/Goods/%s", id)
	return net.GetHtml(url)
}

func extractGoodsInfo(htmlText string) (*GoodsInfo, error) {
	re := regexp.MustCompile(`(?m)<script type="application/ld\+json">((\n|.)*?)</script>`)
	matches := re.FindStringSubmatch(htmlText)
	if len(matches) == 0 {
		return nil, errors.New("goods info not foundin html")
	}

	// parse goods info JSON
	infoJson := matches[1]
	var info GoodsInfo
	if err := json.Unmarshal([]byte(infoJson), &info); err != nil {
		return nil, err
	}

	// unescape title
	info.Name = html.UnescapeString(info.Name)

	return &info, nil
}
