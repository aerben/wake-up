package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

func main() {

	templateFileFlag := flag.String("t", "./resources/template.jpg", "template file path")
	overlayFileFlag := flag.String("l", "./resources/overlay.png", "overlay file path")
	outputFileFlag := flag.String("o", "./output.png", "output file path")
	flag.Parse()

	templateSrc := loadImage(*templateFileFlag)
	overlaySrc := loadImage(*overlayFileFlag)

	preparedOverlay := prepareOverlay(overlaySrc)

	target := image.NewNRGBA(templateSrc.Bounds())

	renderTemplateAndOverlayToTarget(target, templateSrc, preparedOverlay)

	writePng(*outputFileFlag, target)
}

func renderTemplateAndOverlayToTarget(target *image.NRGBA, templateSrc image.Image, preparedOverlay *image.NRGBA) {
	draw.Draw(target, templateSrc.Bounds(), templateSrc, image.Point{X: 0, Y: 0}, draw.Src)
	offset := image.Point{X: -265, Y: -180}
	draw.Draw(target, image.Rect(0, 0, 330, target.Bounds().Dy()), preparedOverlay, offset, draw.Over)
}

func prepareOverlay(overlaySrc image.Image) *image.NRGBA {
	scaledOverlay := resize.Resize(70, 0, overlaySrc, resize.Lanczos3)
	rotated := imaging.Rotate(scaledOverlay, 60, color.Transparent)
	for i := 1; i < 13; i++ {
		drawCircle(rotated, 5, 40, i, color.Transparent)
	}
	return rotated
}

func writePng(path string, target *image.NRGBA) {
	out, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	encodeErr := png.Encode(out, target)
	if encodeErr != nil {
		panic(encodeErr.Error())
	}
}

func loadImage(path string) image.Image {
	infile, err := os.Open(path)
	defer func(infile *os.File) {
		err := infile.Close()
		if err != nil {
			panic(err.Error())
		}
	}(infile)
	if err != nil {
		panic(err.Error())
	}
	src, _, err := image.Decode(infile)
	if err != nil {
		panic(err.Error())
	}
	return src
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
