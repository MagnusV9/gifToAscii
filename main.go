
package gifToAscii

import (
	"image/gif"
	"os"
	"strings"
)

type Pixel struct {
	X, Y int
	R, G, B uint8
}

func getFrameTime(gifPath string) int {
	gif, err := os.Open(gifPath)
	if err != nil {
		panic(err)
	}
	defer gif.Close()
	gifData, err := gif.DecodeAll(gif)
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
				paletteIndex := frame.Pix[j*frame.Stride+i] // frame.Stride er hvor stor en rad i frame arrayet er. Matte for å gå ifra 2 dimensjoner til en.
				color := frame.Palette[paletteIndex]

				r, g, b, _ := color.RGBA()

				pixel := Pixel{
					X: i,
					Y: j,
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
		"@", "@", "@", "%", "%", "#", "*", "*", "+", "+", "=", "-", "-", ":", ".", " ", " ",
	}

	var frames []string

	gifFrames := getGifAsSlice(gifData)

	for _, gifFrame := range gifFrames {
		var builder strings.Builder
		for _, pixel := range gifFrame {
			greyScale := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B) // gjør om fra farge til grå skala.
			builder.WriteString(rgbValueToAscii(greyScale, asciiChars))
		}
		builder.WriteString("\n")
		frames = append(frames, builder.String())
	}

	return frames
}


func clearConsole() {
	switch os := runtime.GOOS; os {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default: // For Linux, MacOS, etc.
		fmt.Print("\033[H\033[2J")
	}
}

func playGif(gifPath string) {
	gifFile, err := os.Open(gifPath)
	if err != nil {
		panic(err)
	}
	defer gifFile.Close()

	gifDecoded, err := gif.DecodeAll(gifFile)
	if err != nil {
		panic(err)
	}

	frames := drawGifFramesToBuffer(gifDecoded)

	for {
		for index, frame := range frames {
			clearConsole()           
      fmt.Println(frame)       
			time.Sleep(time.Duration(gifDecoded.Delay[index]) * 10 * time.Millisecond) // må vente basert på hvor mange fps man skal ha ifra meta dataen til gifen. 
		}
	}
}

func main() {
  if len(os.Args) < 2{
    panic("Need to provide path to gif")
  }

	playGif(os.Args[1])

}
