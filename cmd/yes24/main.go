package main

import (
	"fmt"
	"github.com/librarios/cli/internal/app/yes24"
	"log/slog"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		slog.Error("goods id is not set")
		os.Exit(1)
	}

	id := args[0]

	info, err := yes24.GetBookInfo(id)
	if err != nil {
		slog.Error("failed to get book info", "id", id, "error", err.Error())
		os.Exit(1)
	}

	basename := fmt.Sprintf("%s (%s)", info.Title, info.GetPublishedYear())

	bookInfoFilename := fmt.Sprintf("%s.txt", basename)
	if err = yes24.WriteBookInfo(info, bookInfoFilename); err != nil {
		slog.Error("failed to write book info", "id", id, "error", err.Error())
		os.Exit(1)
	}

	if err := yes24.DownloadPreviewImages(id, basename); err != nil {
		slog.Error("failed to download preview images", "error", err.Error())
		os.Exit(1)
	}
}
