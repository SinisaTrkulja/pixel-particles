package main

import (
	"image/color"
	"math/rand"

	"github.com/gopxl/pixel/v2"
	"golang.org/x/image/colornames"
)

var (
	RED    = Color{to_pixel_color(colornames.Red), 0}
	BLUE   = Color{to_pixel_color(colornames.Blue), 1}
	GREEN  = Color{to_pixel_color(colornames.Green), 2}
	YELLOW = Color{to_pixel_color(colornames.Yellow), 3}
	PURPLE = Color{to_pixel_color(colornames.Purple), 4}
)

func init_particles(particle_count int) ([]Particle, []float32) {
	var colors = []Color{
		RED, BLUE, GREEN, YELLOW, PURPLE,
	}
	particles := make([]Particle, particle_count)
	for i := range particle_count {
		particles[i] = Particle{
			x_position: rand.Float64()*float64(X_MAX_BOUND-X_MIN_BOUND) + X_MIN_BOUND,
			y_position: rand.Float64()*float64(Y_MAX_BOUND-Y_MIN_BOUND) + Y_MIN_BOUND,
			x_speed:    0, // rand.Float64() * SPEED,
			y_speed:    0, // rand.Float64() * SPEED,
			radius:     RADIUS,
			color:      colors[rand.Intn(len(colors))],
			trail:      make([]pixel.Vec, 0, 20),
		}
	}
	// Add positions

	positions_and_velocities := make([]float32, particle_count*5)
	for i := range particles {
		positions_and_velocities[i*5] = float32(particles[i].x_position)
		positions_and_velocities[i*5+1] = float32(particles[i].y_position)
		positions_and_velocities[i*5+2] = float32(particles[i].x_speed)
		positions_and_velocities[i*5+3] = float32(particles[i].y_speed)
		positions_and_velocities[i*5+4] = float32(particles[i].color.position)
	}
	return particles, positions_and_velocities
}

func to_pixel_color(c color.RGBA) pixel.RGBA {
	return pixel.RGBA{
		R: float64(c.R) / 255.0,
		G: float64(c.G) / 255.0,
		B: float64(c.B) / 255.0,
		A: float64(c.A) / 255.0,
	}
}
