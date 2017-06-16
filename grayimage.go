package main


import "os"
import 	"image"
import 	"image/color"
import 	"image/png"


type grayImage struct {
	rect *image.Rectangle
	img *image.Gray 
	f *os.File 
}


func NewGrayImage(r image.Rectangle, fname string) *grayImage {
	
	if thef, err := os.Create(fname) ; err != nil {
		return nil
	} else {
		theimg := image.NewGray(r)
		return &grayImage{rect: &r, f: thef, img: theimg}
	}
}


func (gimg grayImage) Set(x,y, col int)  {
	gray := 255-uint8((col*5) % 256)
//		gray := uint8(255)
	if col == 0 {
		gray = 0
	} 
	gimg.img.SetGray(x,y, color.Gray{gray})
}
func (gimg grayImage) Max() image.Point  { return gimg.rect.Max }
func (gimg grayImage) Min() image.Point  { return gimg.rect.Min }

func (gimg grayImage) 	Sync() error {
	gimg.f.Seek(0,0)
	png.Encode(gimg.f, gimg.img)
	return gimg.f.Sync()
	
}

func (gimg grayImage)  Close() error {
	return gimg.f.Close()
}
