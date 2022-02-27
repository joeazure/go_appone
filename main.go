package main

import (
	"os"

	"jazure.com/go-appone/appone"
)

func main() {
	args := os.Args[1:]
	// srcFile phoneCode align[t, m,  b] outDir
	appone.Wallpaperize(args[0], args[1], args[2], args[3])
}
