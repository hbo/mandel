package main

import "fmt"
import "os"
import "strconv"
import "strings"





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
	case 1:
		if res, err = strconv.Atoi(argsWithoutProg[0]) ; err != nil {
			fmt.Printf("error in parsing parameter %s\n", argsWithoutProg[2])
			panic("panicing!")
		}
	case 2:
		x = tuple_parse(argsWithoutProg[0])
		y = tuple_parse(argsWithoutProg[1])
	case 3:
		x = tuple_parse(argsWithoutProg[1])
		y = tuple_parse(argsWithoutProg[2])
		if res, err = strconv.Atoi(argsWithoutProg[0]) ; err != nil {
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
