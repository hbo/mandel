package main

import "fmt"
//import "io"
import "os"
import "math/cmplx"
import 	"image"
import 	"image/color"
import 	"image/png"

const infinity = 1000
const blocks = 5
const upperbound = 2

//const theX = complex(-1.8625,-0.001)
//const theY = complex(-1.8595,0.001)
const DefaultX = complex(-2.5,-1)
const DefaultY = complex(1.0,1.0)
const DefaultRes = 1000


type imgval struct  {
	x,y, v  int
}


func iter(c complex128, bound float64) int {
        z := complex(0,0)
        zz := complex(0,0)
	n := 0
        for  ; n < infinity ; n++ {
		zz = z*z
                z = zz + c
		abs := cmplx.Abs(z)
                if abs > bound {
//	fmt.Printf("no convergence %C %C %d, %f\n", c, z, n, abs)
                        return n+1
                }
        }
	
//	fmt.Printf("convergence %c %d\n", c, n)
        return 0
}





func mandel(x,y complex128,xpix, ypix, xxpix, yypix int, job int, kanal chan int, imgkanal chan imgval) {

        width := real(y) - real(x)
        height := imag(y) - imag(x)
	pixwidth := xxpix - xpix
	pixheight := yypix - ypix

        if width <= 0 ||
		height <= 0 ||
		pixwidth <= 0 ||
		pixheight <= 0 {
                fmt.Printf("choose your numbers right fool: width %f height %f pixwidth %f pixheight %d\n", width, height, pixwidth, pixheight)
		return
        }

	resx := width / float64(pixwidth)
	resy := height / float64(pixheight)
	for xcount, stepx := xpix,  real(x) ; xcount < xxpix; xcount++ {
		stepx += resx
		for ycount , stepy  := ypix, imag(x) ; ycount < yypix;  ycount++ {
			stepy += resy
                        c := complex(stepx,stepy)
			rounds := iter(c,upperbound)
			imgkanal <- imgval{xcount, ycount, rounds}
                }
        }
	kanal <- job

	
	
	
}

/* func mandel2(x,y complex128, resx, resy float64, job int, kanal chan int, imgkanal chan imgval)  {

        width := real(y) - real(x)
        height := imag(x) - imag(y)
 	fmt.Printf("mandel %d %c %c  width %f height %f \n", job, x, y, width, height )
        if width <= 0 ||  height <= 0 {
                fmt.Printf("choose your numbers fool: width %f height %f \n", width, height)
		return
        }
	xcount := 0
	ycount := 0
        for stepx := real(x) ; stepx < real(y) ; stepx += resx, xcount++ {
                for stepy := imag(x) ; stepy > imag(y) ; stepy -= resy, ycount++ {
                        c := complex(stepx,stepy)
			rounds := iter(c,10e+6)
			imgkanal <- imgval{xcount, stepy, rounds}
                }

        }
	kanal <- job
}
*/

func composeimg(img *image.Gray, kanal chan imgval) {
	points := 0
	for  {
		iv := <-kanal
		points++;
		gray := 255-uint8((iv.v*10) % 256)
//		gray := uint8(255)
		if iv.v == 0 {
			gray = 0
		} 
		x,y  := int(iv.x), int(iv.y)
		
		img.SetGray(x,y, color.Gray{gray})
	}
}

func parse_args() (complex128, complex128, int) {

	res := DefaultRes
	x := DefaultX
	y := DefaultY
	

	
	return x, y, res
	
}

func main() {

	theX, theY, res := parse_args()

        width := real(theY) - real(theX)
        height := imag(theY) - imag(theX)

	ratio := width/height
	
	pixx := res
	pixy := int(float64(res)/ratio)

	fmt.Printf("generating pic with ratio %f, %dx%d\n", ratio, pixx, pixy)

	xcornerx := -pixx/2
	xcornery := -pixy/2
	ycornerx := pixx/2
	ycornery := pixy/2
	
	rect := image.Rect( xcornerx, xcornery, ycornerx, ycornery )
	
	img := image.NewGray(rect)
	
	
	
	sizex := width / blocks
	sizey := height / blocks

	blocksizex  := pixx / blocks
	blocksizey  := pixy / blocks
	
	job := 0
	kanal :=  make(chan int)
	imgkanal :=  make(chan imgval)

	go composeimg(img, imgkanal)

	for i := 0 ; i < blocks; i++ {
		dx :=  real(theX) + float64(i) * sizex
		ddx :=  real(theX) + float64(i+1) * sizex

		leftx := xcornerx + i * blocksizex 
		rightx := xcornerx + (i+1) * blocksizex 

		for j := 0 ; j < blocks; j++ {
			dy :=  imag(theX) + float64(j) * sizey
			ddy :=  imag(theX) + float64(j+1) * sizey

			lowery := xcornery + j * blocksizey
			uppery := xcornery + (j+1) * blocksizey
			
			go mandel(complex(dx, dy), complex(ddx, ddy), leftx, lowery, rightx,uppery, job, kanal, imgkanal)
			job++
		}
	}

	f , _ := os.Create("/tmp/img.png")
	for ; job > 0 ; job-- {
		f.Seek(0,0)
		png.Encode(f, img)
		f.Sync()
		k:= <-kanal
		fmt.Printf("done %d\n", k)
	}
	f.Close()

//	fmt.Printf("img minx %d miny %d maxx %d maxy %d \n",  img.Bounds().Min.X, img.Bounds().Min.Y, img.Bounds().Max.X, img.Bounds().Max.Y  )


	
}
