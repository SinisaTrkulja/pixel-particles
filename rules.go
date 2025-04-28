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

func interaction(pa, pb Particle) (float64, float64) {
	return interactionMatrix[pa.color.position][pb.color.position],
		interactionMatrix[pb.color.position][pb.color.position]
}

func key_listener(win *opengl.Window, particles *[]Particle, positions_and_velocities *[]float32) {
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
		*particles, *positions_and_velocities = init_particles(PARTICLE_COUNT)
	case win.JustPressed(pixel.KeyUp):
		new_particles, new_positions_and_velocities := init_particles(COUNT_CHANGE_STEP)
		*particles = append(*particles, new_particles...)
		*positions_and_velocities = append(*positions_and_velocities, new_positions_and_velocities...)
		PARTICLE_COUNT += COUNT_CHANGE_STEP
	case win.JustPressed(pixel.KeyDown):
		*particles = (*particles)[COUNT_CHANGE_STEP:]
		*positions_and_velocities = (*positions_and_velocities)[COUNT_CHANGE_STEP*5:]
		PARTICLE_COUNT -= COUNT_CHANGE_STEP
	}
}
