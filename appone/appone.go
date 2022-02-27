package appone

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

var phoneSizes = map[string][]int{
	"i_xr":   {828, 1792},
	"i_xsx":  {1242, 2688},
	"i_xs":   {1125, 2436},
	"i_x":    {1125, 2436},
	"i_11px": {1242, 2688},
	"i_11p":  {1125, 2436},
	"i_11":   {828, 1792},
	"i_12pm": {1284, 2778},
	"i_12p":  {1170, 2532},
	"i_12m":  {1125, 2436},
	"i_12":   {1170, 2532},
	"i_13pm": {1284, 2778},
	"i_13p":  {1170, 2532},
	"i_13m":  {1080, 2340},
	"i_13":   {1170, 2532},
}

func getDimensions(code string) []int {
	return phoneSizes[code]
}

func isLandscape(w int, h int) bool {
	return w >= h
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func createNewFilledWallpaper(w int, h int, bg color.Color) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{w, h}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	return img
}

func Wallpaperize(srcFileName string, phoneCode string, align string, outDir string) {
	// Read image from file that already exists
	existingImageFile, err := os.Open(srcFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer existingImageFile.Close()

	// Since we know it is a png already can call png.Decode() directly
	// Otherwise use generic image.Decode() and look at the type if necessary
	loadedImage, err := png.Decode(existingImageFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image loaded and decoded")
	originalSize := loadedImage.Bounds().Max
	// Find out the background color (color at pixel 0,0)
	pixColor := loadedImage.At(0, 0)
	fmt.Println("Color at 0,0 is", pixColor)

	wallpaperSize := getDimensions(phoneCode)
	fmt.Println(wallpaperSize)
	wallImage := createNewFilledWallpaper(wallpaperSize[0], wallpaperSize[1], pixColor)

	// Scale original down based on orientation
	if isLandscape(originalSize.X, originalSize.Y) {
		fmt.Println("Landscape")
		dstImage := imaging.Resize(loadedImage, wallpaperSize[0], 0, imaging.Lanczos)

		var y int
		if align == "t" {
			//top
			y = 0
		} else if align == "m" {
			// middle
			y = (wallpaperSize[1] - dstImage.Rect.Size().Y) / 2
		} else {
			// bottom
			y = (wallpaperSize[1] - dstImage.Rect.Size().Y)
		}
		offset := image.Point{0, y}
		// Draw the scaled image onto wallpaper
		draw.Draw(wallImage, dstImage.Bounds().Add(offset), dstImage, image.Point{}, draw.Over)
	} else {
		fmt.Println("Portrait")
		// TODO
	}
	fName := fmt.Sprint(fileNameWithoutExt(srcFileName), "-", phoneCode, "-", align, ".png")

	_, err = os.Stat(outDir)
	if os.IsNotExist(err) {
		fmt.Println("Creating output directory: ", outDir)
		err = os.Mkdir(outDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Full Path
	fName = filepath.Join(outDir, fName)
	f, _ := os.Create(fName)
	png.Encode(f, wallImage)
}
