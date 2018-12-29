package AnimaKit

import (
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
