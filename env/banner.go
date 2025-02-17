package env

import (
	"embed"
	"fmt"
)

//go:embed banner.txt
var banner embed.FS

func printBanner() {
	if enablePrintBanner := envMap[AliothFrameworkBannerPrintKey]; enablePrintBanner == "TRUE" {
		bannerBytes, err := banner.ReadFile("banner.txt")
		if err == nil {
			// try to print banner, if error occurs, ignore it
			// you can delete this block if you don't want to print banner,
			// or you can change the banner.txt file to customize your banner
			fmt.Println(string(bannerBytes))
		}
	}
}
