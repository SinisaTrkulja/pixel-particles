package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

const (
	PARTICLE_COUNT = 6000
	FORCE_RANGE    = 60.0
	X_MIN          = 0.0
	X_MAX          = 1024.0
	Y_MIN          = 0.0
	Y_MAX          = 700
	RADIUS         = 2.0
	CHARGE         = 1.7
	DELTA          = 1.0
	PROXIMA_DAMP   = 1.0
	DISTAL_DAMP    = 0.9995
	EPSILON        = 0.2

	X_MIN_BOUND = X_MIN + RADIUS + WALL_OFFSET
	X_MAX_BOUND = X_MAX - RADIUS - WALL_OFFSET
	Y_MIN_BOUND = Y_MIN + RADIUS + WALL_OFFSET
	Y_MAX_BOUND = Y_MAX - RADIUS - WALL_OFFSET
	WALL_OFFSET = 55.0
)

func main() {
	opengl.Run(run)
}

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Particle Simulator",
		Bounds: pixel.R(X_MIN, Y_MIN, X_MAX, Y_MAX),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	particles := init_particles(PARTICLE_COUNT)

	for !win.Closed() {
		if win.JustPressed(pixel.KeyP) {
			PAUSED = !PAUSED // Toggle pause state
		}
		win.Clear(colornames.Black)
		if !PAUSED {
			update_particles(particles)
		}
		draw_particles(win, particles)
		win.Update()
	}
}

func update_particles(particles []Particle) {
	update_velocities(particles)
	update_positions(particles)
}

func draw_particles(win *opengl.Window, particles []Particle) {
	imd := imdraw.New(nil)
	for _, particle := range particles {
		imd.Color = particle.color.rgba
		imd.Push(pixel.V(particle.x_position, particle.y_position))
		imd.Circle(particle.radius, 0)
	}
	imd.Draw(win)
}
