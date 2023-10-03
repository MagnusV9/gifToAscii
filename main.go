package gifToAscii

import (
	"os"
  "image/gif"
  "fmt"
)

func getFrameTime(gifPath string) int {
	gif, err := os.Open(gifPath)
	if err != nil {
		panic(err)
	}
	gifData, err := gif.DecodeAll(gif)
	if err != nil {
		panic(err)
	}
	return gifData.Delay
}

func getGifAsSlice(gif image.Gif){
  var frames [][]Pixel 

  for _, frame := range gifData.Image{
    var framePixels []Pixel

    bounds := frame.Bounds()

    for i := _, bounds.Min.Y; i < bounds.Max.Y; i ++ {
      for j := _ , bounds.Min.X; j < bounds.Max.X; j++{
        color := frameAt(i,j)
        
        r, g, b := color.RGBA()

        pixel := Pixel{
          X: x,
          Y: y, 
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

func drawGifFramesToBuffer(gif image.Gif){ // tegner først representasjonen å legger det inn i en slice for å hindre mange io operasjoner da disse er trege.
  var asciiChars = []string{
	  "@", "@", "@", "%", "%", "#", "*", "*", "+", "+", "=", "-", "-", ":", ".", " ", " ",
  }

  var frames [] string

  gifFrames := getGifAsSlice(gif)


    for i := 0; i < length(gifFrames){
      for j := 0; j < length(gifFrames[0]){


      }
    }

}

func main() {

	gifData, err := os.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

}
