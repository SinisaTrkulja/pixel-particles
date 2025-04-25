package main

import "github.com/gopxl/pixel/v2"

type Color struct {
	rgba     pixel.RGBA
	position int
}

type Particle struct {
	x_position, y_position, x_speed, y_speed float64
	radius                                   float64
	color                                    Color
}
