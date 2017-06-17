package main

import 	"image"

type TargetImage interface {
	Set(x,y, col int)
	Max() image.Point
	Min() image.Point
	Sync() error
	Close() error
}

