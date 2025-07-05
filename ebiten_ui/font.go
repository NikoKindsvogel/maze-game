package ebiten_ui

import (
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var MainFont font.Face
var HeadlineFont font.Face

func init() {
	fontBytes, err := ioutil.ReadFile("assets/fonts/times.ttf")
	if err != nil {
		log.Fatalf("Failed to load font file: %v", err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}

	MainFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Failed to create font face: %v", err)
	}

	HeadlineFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    headlineSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Failed to create font face: %v", err)
	}
}
