package main

import "fmt"
//import "io"
import "os"
import "math"
import "math/cmplx"
import 	"image"
import "time"


const infinity = 1000
const blocks = 5
const upperbound = 2

//const theX = complex(-1.8625,-0.001)
//const theY = complex(-1.8595,0.001)
const DefaultX = complex(-2.1,-1.2)
const DefaultY = complex(0.5,1.2)
const DefaultRes = 1000


type imgval struct  {
	x,y, v  int
}

type TargetImage interface {
	Set(x,y, col int)
	Max() image.Point
	Min() image.Point
	Sync() error
	Close() error
}




func composeimg(img TargetImage, kanal chan imgval, quit chan bool) {
	
	
	nrpoints := (img.Max().X - img.Min().X)  * (img.Max().Y - img.Min().Y)
	writepoints := int(nrpoints/10)
	
	//	img := image.NewGray(rect)

	
	defer func (quit chan  bool ){
		img.Sync()
		img.Close()
		fmt.Printf("composeimg: quitting\n")
		quit <- true
	}(quit)
	
	points := 0
	var iv imgval
	do_quit := false
	for  {
		select {
		case iv = <-kanal:
		case <-quit:
			do_quit = true
			fmt.Printf("composeimg: received quit\n")
		case <-time.After(time.Second * 1):
			if do_quit { return }
		}
		
		
		x,y  := int(iv.x), int(iv.y)
		
		img.Set(x,y, iv.v)

		if points % writepoints == 0 {
			img.Sync()
		}
		points++;
	}
	
	
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
                        return n+1
                }
        }
	
        return 0
}



func trivial_point ( x , y float64  ) bool {
	xless := (x - 0.25)

	yy := y*y
	pp := xless*xless + yy
	p := math.Sqrt( pp )
	return ( x < p - 2*pp + 0.25  ) || ( (x+1)*(x+1) + yy < 1/16.0 )
	
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
			var rounds int
			if trivial_point(stepx,stepy) {
				rounds = 0
			} else {
				rounds = iter(c,upperbound)
			}
			imgkanal <- imgval{xcount, ycount, rounds}
                }
        }
	kanal <- job

	
	
	
}



func main() {

	defer func () {
		if err := recover(); err != nil {
			fmt.Printf("too bad, we have to say: %s\n", err)
			return
		}
	}()
	
	theX, theY, res := parse_args()


        width := real(theY) - real(theX)
        height := imag(theY) - imag(theX)

	ratio := width/height
	
	pixx := res
	pixy := int(float64(res)/ratio)

	fmt.Printf("generating pic with ratio %f, %dx%d\n", ratio, pixx, pixy)
	fmt.Printf("%s %d %f,%f %f,%f\n", os.Args[0], res, real(theX), imag(theX), real(theY), imag(theY))
	

	xcornerx := -pixx/2
	xcornery := -pixy/2
	ycornerx := pixx/2
	ycornery := pixy/2
	
	rect := image.Rect( xcornerx, xcornery, ycornerx, ycornery )
	
	
//	img := NewGrayImage( rect, "/tmp/gray.png")
	img := NewColImage( rect, "/tmp/col.png")
	
	sizex := width / blocks
	sizey := height / blocks

	blocksizex  := pixx / blocks
	blocksizey  := pixy / blocks
	
	job := 0
	kanal :=  make(chan int)
	imgkanal :=  make(chan imgval)
	quit :=  make(chan bool)
	httpquit :=  make(chan bool)

	go composeimg(img, imgkanal, quit)

	go httpserver(img, httpquit)

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

	for ; job > 0 ; job-- {
		k:= <-kanal
		fmt.Printf("done %d\n", k)
	}

	// signale composeimg goroutine that all threads are dead
	quit <- true
	// wait for ack from composeimg
	<-quit  

	<-httpquit
//	fmt.Printf("img minx %d miny %d maxx %d maxy %d \n",  img.Bounds().Min.X, img.Bounds().Min.Y, img.Bounds().Max.X, img.Bounds().Max.Y  )


	
}
