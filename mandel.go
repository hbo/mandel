package main

import "fmt"
import "math/cmplx"


const infinity = 10000
const blocks = 2


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

func mandel(x,y complex128, resx, resy float64, job int, kanal chan int)  {

        width := real(y) - real(x)
        height := imag(x) - imag(y)
 	fmt.Printf("mandel %d %c %c  width %f height %f \n", job, x, y, width, height )
        if width <= 0 ||  height <= 0 {
                fmt.Printf("choose your numbers fool: width %f height %f \n", width, height)
		return
        }
        for stepx := real(x) ; stepx < real(y) ; stepx += resx {
                for stepy := imag(x) ; stepy > imag(y) ; stepy -= resy {
                        c := complex(stepx,stepy)
			iter(c,10e+6)
                }

        }
	kanal <- job
}

func main() {


	x := complex(-2,1)
	y := complex(2,-1)

        width := real(y) - real(x)
        height := imag(x) - imag(y)

	pixx := 1000
	pixy := 500
	
	resx := width / float64(pixx)
	resy := height / float64(pixy)
	
	sizex := width / blocks
	sizey := height / blocks
	job := 0
	kanal :=  make(chan int)
	for i := 0 ; i < blocks; i++ {
		dx :=  real(x) + float64(i) * sizex
		ddx :=  real(x) + float64(i+1) * sizex
		for j := blocks ; j > 0; j-- {
			dy :=  imag(y) + float64(j) * sizey
			ddy :=  imag(y) + float64(j-1) * sizey
			go mandel(complex(dx, dy), complex(ddx, ddy), resx, resy, job, kanal)
			job++
		}
	}

	for ; job > 0 ; job-- {
		k:= <-kanal
		fmt.Printf("done %d\n", k)
	}
	
}
