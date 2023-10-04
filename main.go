package main

import (
  "image/draw"
	"image/gif"
	"os"
	"strings"
	"fmt"
	"time"
  "github.com/nfnt/resize"
  "image"
)


type Pixel struct {
	X, Y int
	R, G, B uint8
}

func getFrameTime(gifPath string) int {
	gifFile, err := os.Open(gifPath)
	if err != nil {
		panic(err)
	}
	defer gifFile.Close()
  gifData, err := gif.DecodeAll(gifFile)
	if err != nil {
		panic(err)
	}
	return gifData.Delay[0]
}


func getGifAsSlice(gifData *gif.GIF) [][]Pixel {
	var frames [][]Pixel

	for _, frame := range gifData.Image {

		var framePixels []Pixel

		bounds := frame.Bounds()

		for i := bounds.Min.Y; i < bounds.Max.Y; i++ {
			for j := bounds.Min.X; j < bounds.Max.X; j++ {
				paletteIndex := frame.Pix[i*frame.Stride+j] // frame.Stride er hvor stor en rad i frame arrayet er. Matte for å gå ifra 2 dimensjoner til en.
				color := frame.Palette[paletteIndex]

				r, g, b, _ := color.RGBA()

				pixel := Pixel{
					X: j,
					Y: i,
					R: uint8(r >> 8), // Siden pixelene er representert som hex kan vi høyre skifte binærtallet slik at vi får en 256 bit representasjon istedet.
					G: uint8(g >> 8),
					B: uint8(b >> 8),
				}
				framePixels = append(framePixels, pixel)
			}
		}
		frames = append(frames, framePixels)
	}

	return frames
}

func rgbValueToAscii(intensity float64, asciiChars []string) string {
	asciiIndex := int(intensity*float64(len(asciiChars))/256) % len(asciiChars)
	return asciiChars[asciiIndex]
}

func drawGifFramesToBuffer(gifData *gif.GIF) []string {
	asciiChars := []string{
    "@", "%", "#", "*", "+", "=", "-", ":", ".", " ",
	}

	var frames []string

	gifFrames := getGifAsSlice(gifData)
  for _, gifFrame := range gifFrames {
    var builder strings.Builder
    width := gifData.Image[0].Bounds().Dx()
    counter := 0
    for _, pixel := range gifFrame {
        greyScale := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)
        builder.WriteString(rgbValueToAscii(greyScale, asciiChars))
        counter++
        if counter == width {
            builder.WriteString("\n")
            counter = 0
        }
    }
    frames = append(frames, builder.String())
}

	return frames
}


func clearConsole() {
	// For Linux, MacOS, windows idk.
	fmt.Print("\033[H\033[2J")
}


func playGif(gifPath string, newWidth, newHeight uint) {
    gifFile, err := os.Open(gifPath)
    if err != nil {
        panic(err)
    }
    defer gifFile.Close()

    gifDecoded, err := gif.DecodeAll(gifFile)
    if err != nil {
        panic(err)
    }

    // Resize each frame of the gif
    for i, frame := range gifDecoded.Image {
        resizedImage := resize.Resize(newWidth, newHeight, frame, resize.Lanczos3)
        palettedImage := image.NewPaletted(resizedImage.Bounds(), frame.Palette)
        draw.Draw(palettedImage, palettedImage.Rect, resizedImage, resizedImage.Bounds().Min, draw.Over)
        gifDecoded.Image[i] = palettedImage

            }

    frames := drawGifFramesToBuffer(gifDecoded)

    for {
        for index, frame := range frames {
            clearConsole()
            fmt.Println(frame)
            time.Sleep(time.Duration(gifDecoded.Delay[index]) * 10 * time.Millisecond)
        }
    }
}

func main() {
  if len(os.Args) < 2{
    panic("Need to provide path to gif")
  }

	playGif(os.Args[1], 300, 100)

}
