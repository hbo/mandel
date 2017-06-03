package main

import "fmt"
//import "io"
import "os"
import "math"
import "math/cmplx"
import 	"image"
import 	"image/color"
import 	"image/png"
import "strconv"
import "strings"
import "time"


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

func grayimg(rect image.Rectangle, kanal chan imgval, quit chan bool) {
	

	nrpoints := (rect.Max.X - rect.Min.X)  * (rect.Max.Y - rect.Min.Y)
	writepoints := int(nrpoints/10)
	
	img := image.NewGray(rect)

	f , _ := os.Create("/tmp/gray.png")

	defer func (f *os.File, quit chan  bool ){
		f.Seek(0,0)
		png.Encode(f, img)
		f.Sync()
		f.Close()
		fmt.Printf("grayimg: quitting\n")
		quit <- true
	}(f,quit)
	
	points := 0
	var iv imgval
	do_quit := false
	for  {
		select {
		case iv = <-kanal:
		case <-quit:
			do_quit = true
			fmt.Printf("grayimg: received quit\n")
		case <-time.After(time.Second * 1):
			if do_quit { return }
		}
		
		
		gray := 255-uint8((iv.v*10) % 256)
//		gray := uint8(255)
		if iv.v == 0 {
			gray = 0
		} 
		x,y  := int(iv.x), int(iv.y)
		
		img.SetGray(x,y, color.Gray{gray})

		if points % writepoints == 0 {
			f.Seek(0,0)
			png.Encode(f, img)
			f.Sync()
		}
		points++;
	}
	
	
}

func tuple_parse ( tuple string ) complex128 {

	comps := strings.Split(tuple, ",")

	real , err  := strconv.ParseFloat(comps[0], 64)

	if err != nil {
		fmt.Printf("error in parsing parameter %s\n", tuple)
		panic("panicing!")
	}
	
	img , err  := strconv.ParseFloat(comps[1], 64)

	if err != nil {
		fmt.Printf("error in parsing parameter %s\n", tuple)
		panic("panicing!")
	}

	return complex(real, img)
}


func parse_args() (complex128, complex128, int) {

	var err error

	res := DefaultRes
	x := DefaultX
	y := DefaultY


	// We expect either 0 arguments (in which case the defaults hold)or 2  (x,y coords), or 3 (x,y,res)
	
	argsWithoutProg := os.Args[1:]

	switch len(argsWithoutProg) {
	case 2:
		x = tuple_parse(argsWithoutProg[0])
		y = tuple_parse(argsWithoutProg[1])
	case 3:
		x = tuple_parse(argsWithoutProg[0])
		y = tuple_parse(argsWithoutProg[1])
		if res, err = strconv.Atoi(argsWithoutProg[2]) ; err != nil {
			fmt.Printf("error in parsing parameter %s\n", argsWithoutProg[2])
			panic("panicing!")
		}
	case 0:
		return x, y, res
	default:
		fmt.Printf("error in parsing parameter: %d parameters\nshould be 0, 2 or 3\n", len(argsWithoutProg))
		panic("panicing!")
	}
	
	return x, y, res
	
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

	xcornerx := -pixx/2
	xcornery := -pixy/2
	ycornerx := pixx/2
	ycornery := pixy/2
	
	rect := image.Rect( xcornerx, xcornery, ycornerx, ycornery )
	
//	img := image.NewGray(rect)
	
	
	
	sizex := width / blocks
	sizey := height / blocks

	blocksizex  := pixx / blocks
	blocksizey  := pixy / blocks
	
	job := 0
	kanal :=  make(chan int)
	imgkanal :=  make(chan imgval)
	quit :=  make(chan bool)

	go grayimg(rect, imgkanal, quit)

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

	quit <- true

	_ = <-quit  

	

	
	
//	fmt.Printf("img minx %d miny %d maxx %d maxy %d \n",  img.Bounds().Min.X, img.Bounds().Min.Y, img.Bounds().Max.X, img.Bounds().Max.Y  )


	
}
