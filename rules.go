package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

var interactionMatrix = [][]float64{
	//  Red,   Blue,  Green, Yellow, Purple
	{0.04, 0.05, -0.02, -0.03, -0.01}, // Red
	{-0.04, -0.01, -0.03, 0.04, 0.02}, // Blue
	{0.03, -0.03, -0.07, -0.04, 0.01}, // Green
	{-0.02, 0.03, -0.03, 0.02, 0.04},  // Yellow
	{0.01, 0.02, 0.01, 0.03, 0.015},   // Purple
}

func interaction(acted_upon, acting Particle) float64 {
	return interactionMatrix[acted_upon.color.position][acting.color.position]
}

func key_listener(win *opengl.Window, particles *[]Particle) {
	switch {
	case win.JustPressed(pixel.KeyP):
		PAUSED = !PAUSED
	case win.JustPressed(pixel.KeyF):
		FADE_TRAIL = !FADE_TRAIL
	case win.JustPressed(pixel.KeyC):
		CHARGE += 0.1
	case win.JustPressed(pixel.KeyV):
		CHARGE -= 0.1
	case win.JustPressed(pixel.KeyY):
		EPSILON += 0.1
	case win.JustPressed(pixel.KeyV):
		EPSILON -= 0.1
	case win.JustPressed(pixel.KeyR):
		*particles = init_particles(PARTICLE_COUNT)
	case win.JustPressed(pixel.KeyUp):
		*particles = append(*particles, init_particles(COUNT_CHANGE_STEP)...)
	case win.JustPressed(pixel.KeyDown):
		*particles = (*particles)[COUNT_CHANGE_STEP:]
	}

}
