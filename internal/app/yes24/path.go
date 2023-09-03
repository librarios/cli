package yes24

import "strings"

func NormalizeFilename(filename string) string {
	replaceMap := map[string]string{
		"?":  "？",
		":":  "：",
		"/":  "／",
		"\\": "＼",
		"*":  "×",
		"\"": "＂",
		"<":  "＜",
		">":  "＞",
		"|":  "｜",
	}

	result := filename
	for from, to := range replaceMap {
		result = strings.ReplaceAll(result, from, to)
	}
	return result
}
