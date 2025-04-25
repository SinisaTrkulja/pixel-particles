package main

import (
	"sync"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

const (
	PARTICLE_COUNT = 4000
	FORCE_RANGE    = 60.0
	X_MIN          = 0.0
	X_MAX          = 1024.0
	Y_MIN          = 0.0
	Y_MAX          = 700
	RADIUS         = 3.0
	CHARGE         = 1.7
	DELTA          = 1.0
	PROXIMA_DAMP   = 1.0
	DISTAL_DAMP    = 0.9995
	EPSILON        = 0.2
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
		win.Clear(colornames.Black)
		update_particles(win, particles)
		win.Update()
	}
}

func update_particles(win *opengl.Window, particles []Particle) {
	calcute_velocities(particles)
	draw_particles(win, particles)
}

func draw_particles(win *opengl.Window, particles []Particle) {
	var wg sync.WaitGroup
	for i := range particles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			particle := &particles[i]
			update_positions(particle)

			imd := imdraw.New(nil)
			imd.Color = particle.color.rgba
			imd.Push(pixel.V(particle.x_position, particle.y_position))
			imd.Circle(particle.radius, 0)
			imd.Draw(win)
		}()
	}
	wg.Wait()

}
