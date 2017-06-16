package main

import "os"
import 	"image"
import 	"image/color"
import 	"image/png"



type colImage struct {
	rect *image.Rectangle
	img *image.NRGBA 
	f *os.File 
}



func NewColImage(r image.Rectangle, fname string) *colImage {
	
	if thef, err := os.Create(fname) ; err != nil {
		return nil
	} else {
		theimg := image.NewNRGBA(r)
		return &colImage{rect: &r, f: thef, img: theimg}
	}
}






func (gimg colImage)  Close() error {
	return gimg.f.Close()
}

func (gimg colImage) 	Sync() error {
	gimg.f.Seek(0,0)
	png.Encode(gimg.f, gimg.img)
	return gimg.f.Sync()
	
}

func (gimg colImage) Max() image.Point  { return gimg.rect.Max }
func (gimg colImage) Min() image.Point  { return gimg.rect.Min }

func (gimg colImage) Set(x,y, col int)  {
	r, g, b := uint8(0), uint8(0), uint8(0)
	

	if col != 0 {
		// this here defines the coloring of the picture
		// currently quite daftly made.  We probably should
		// define a curve through rgb space.
		r = uint8((7*col + 25) % 256 )
		b = uint8((11*col + 50) % 256)
		g = uint8((17*col + 75) % 256)
	}
	gimg.img.Set(x,y, color.NRGBA{r,g,b,255})
}

