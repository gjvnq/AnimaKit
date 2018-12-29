package AnimaKit

import (
	"fmt"
	"image/gif"
	"os"
	"strconv"

	sdlImg "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	IMG_UNKNOWN = 0
	IMG_BMP     = 1
	IMG_CUR     = 2
	IMG_GIF     = 3
	IMG_ICO     = 4
	IMG_JPG     = 5
	IMG_LBM     = 6
	IMG_PCX     = 7
	IMG_PNG     = 8
	IMG_PNM     = 9
	IMG_TIF     = 10
	IMG_WEBP    = 11
	IMG_XCF     = 12
	IMG_XV      = 13
)

func getImgFormat(src *sdl.RWops) int {
	switch {
	case sdlImg.IsBMP(src):
		return IMG_BMP
	case sdlImg.IsCUR(src):
		return IMG_CUR
	case sdlImg.IsGIF(src):
		return IMG_GIF
	case sdlImg.IsICO(src):
		return IMG_ICO
	case sdlImg.IsJPG(src):
		return IMG_JPG
	case sdlImg.IsLBM(src):
		return IMG_LBM
	case sdlImg.IsPCX(src):
		return IMG_PCX
	case sdlImg.IsPNG(src):
		return IMG_PNG
	case sdlImg.IsPNM(src):
		return IMG_PNM
	case sdlImg.IsTIF(src):
		return IMG_TIF
	case sdlImg.IsWEBP(src):
		return IMG_WEBP
	case sdlImg.IsXCF(src):
		return IMG_XCF
	case sdlImg.IsXCF(src):
		return IMG_XCF
	case sdlImg.IsXV(src):
		return IMG_XV
	default:
		return IMG_UNKNOWN
	}
}

type Image struct {
	Frames    []*sdl.Texture
	LenFrames int
	Segs      []CyclableSegments
}

func loadImageFromFile(path string) Image {
	file, err := os.Open(path)
	panicOnError(err)

	rwops := sdl.RWFromFile(path, "r")
	format := getImgFormat(rwops)
	fmt.Println(format)

	if format != IMG_GIF {
		panic("unimplemented format: " + strconv.Itoa(format))
	}

	my_gif, err := gif.DecodeAll(file)
	panicOnError(err)

	fmt.Println(my_gif)

	return Image{}
}

func (self Image) Frame(frame int) *sdl.Texture {
	// Static images
	if self.LenFrames == 0 {
		return nil
	}
	if self.LenFrames == 1 || len(self.Segs) == 0 {
		return self.Frames[0]
	}

	// Animated stuff
	for _, seg := range self.Segs {
		if seg.StartFrame <= frame && frame < seg.StopFrame {
			return self.Frames[seg.WhichFrame(frame)]
		}
	}
	return nil
}

type CyclableSegments struct {
	StartFrame int
	StopFrame  int
	// Both trims are inclusive
	TrimStart int
	TrimEnd   int
	Speed     float64
	Cycle     bool
}

func (self CyclableSegments) WhichFrame(current_frame int) int {
	current_frame -= self.StartFrame

	frame_walk := float64(current_frame) * self.Speed
	out_frame := current_frame + self.TrimStart + int(frame_walk)
	self_len := self.TrimEnd - self.TrimStart
	if !self.Cycle {
		return out_frame
	}
	for out_frame > self.TrimEnd {
		out_frame -= self_len
	}
	for out_frame < self.TrimStart {
		out_frame += self_len
	}
	return out_frame
}
