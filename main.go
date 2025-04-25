package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

const (
	FORCE_RANGE = 60.0
	X_MIN       = 0.0
	X_MAX       = 1424.0
	Y_MIN       = 0.0
	Y_MAX       = 800.0
	RADIUS      = 3.0

	DELTA         = 1.0
	PROXIMAL_DAMP = 1.0
	DISTAL_DAMP   = 0.9995

	X_MIN_BOUND       = X_MIN + RADIUS + WALL_OFFSET
	X_MAX_BOUND       = X_MAX - RADIUS - WALL_OFFSET
	Y_MIN_BOUND       = Y_MIN + RADIUS + WALL_OFFSET
	Y_MAX_BOUND       = Y_MAX - RADIUS - WALL_OFFSET
	WALL_OFFSET       = 25.0
	COUNT_CHANGE_STEP = 100
)

var (
	PARTICLE_COUNT = 4000
	CHARGE         = 1.7
	EPSILON        = 0.2
	PAUSED         = false
	FADE_TRAIL     = false
)

func main() {
	opengl.Run(run)
}

func run() {
	cfg := opengl.WindowConfig{
		Title:                  "Particle Simulator",
		Bounds:                 pixel.R(X_MIN, Y_MIN, X_MAX, Y_MAX),
		VSync:                  true,
		TransparentFramebuffer: true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	particles := init_particles(PARTICLE_COUNT)

	//win.SetSmooth(true)
	for !win.Closed() {
		key_listener(win, &particles)
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
	for i := range particles {
		particle := &particles[i]
		if FADE_TRAIL {
			trail_animation(particle, imd)
		}

		imd.Color = particle.color.rgba
		imd.Push(pixel.V(particle.x_position, particle.y_position))
		imd.Circle(particle.radius, 0)
	}
	imd.Draw(win)
}

func trail_animation(particle *Particle, imd *imdraw.IMDraw) {
	particle.trail = append(particle.trail, pixel.V(particle.x_position, particle.y_position))
	if len(particle.trail) > 4 {
		particle.trail = particle.trail[1:]
	}
	for i, pos := range particle.trail {
		t := float64(i) / float64(len(particle.trail))
		alpha := t * 0.5
		rgba := particle.color.rgba
		rgba.A = alpha
		imd.Color = rgba
		imd.Push(pos)
		radius := particle.radius * (1.0 + 0.5*(1-t))
		imd.Circle(radius, 0)
	}
}
