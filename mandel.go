package main

import "fmt"
import "math/cmplx"
import 	"image"
import 	"image/color"

const infinity = 10000
const blocks = 2

type imgval struct  {
	x,y float64
	v   int
}


func iter(c complex128, bound float64) int {
        z := complex(0,0)
        for n := 0  ; n < infinity ; n++ {
                z := z*z + c
                if cmplx.Abs(z) > bound {
                        return n+1
                }
        }
        return 0
}

func mandel(x,y complex128, resx, resy float64, job int, kanal chan int, imgkanal chan imgval)  {

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


func composeimg(img *image.Gray16, kanal chan imgval) {
	for true {
		iv := <-kanal
		gray := uint8(iv.v % 256)
		x,y  := int(iv.x), int(iv.y)
		fmt.Printf("composeimg %d %d %d    %f   %f  \n",  x, y, gray, iv.x, iv.y )
		
		img.Set(x,y, color.Gray{gray})
	}
}

func main() {


	x := complex(-2,1)
	y := complex(2,-1)

        width := real(y) - real(x)
        height := imag(x) - imag(y)

	pixx := 1000
	pixy := 500

	rect := image.Rect( -pixx/2, pixy/2, pixx/2, -pixy/2)
	img := image.NewGray16(rect)
	
	
	resx := width / float64(pixx)
	resy := height / float64(pixy)
	
	sizex := width / blocks
	sizey := height / blocks
	job := 0
	kanal :=  make(chan int)
	imgkanal :=  make(chan imgval)

	go composeimg(img, imgkanal)

	for i := 0 ; i < blocks; i++ {
		dx :=  real(x) + float64(i) * sizex
		ddx :=  real(x) + float64(i+1) * sizex
		for j := blocks ; j > 0; j-- {
			dy :=  imag(y) + float64(j) * sizey
			ddy :=  imag(y) + float64(j-1) * sizey
			go mandel(complex(dx, dy), complex(ddx, ddy), resx, resy, job, kanal, imgkanal)
			job++
		}
	}

	for ; job > 0 ; job-- {
		k:= <-kanal
		fmt.Printf("done %d\n", k)
	}
	
}
