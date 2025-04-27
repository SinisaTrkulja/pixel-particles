package main

import (
	"fmt"
	"time"

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
	PARTICLE_COUNT int = 1024
	CHARGE             = 1.7
	EPSILON            = 0.2
	PAUSED             = true
	FADE_TRAIL         = false
	RAN                = false

	positions_and_velocities = make([]float32, PARTICLE_COUNT*5)

	particles = init_particles(PARTICLE_COUNT)

	interactionMatrix2 = []float32{
		0.04, 0.05, -0.02, -0.03, -0.01, // Red
		-0.04, -0.01, -0.03, 0.04, 0.02, // Blue
		0.03, -0.03, -0.07, -0.04, 0.01, // Green
		-0.02, 0.03, -0.03, 0.02, 0.04, // Yellow
		0.01, 0.02, 0.01, 0.03, 0.015, // Purple
	}
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

	context, queue, kernel := init_kernel()

	//win.SetSmooth(true)
	var last_time time.Time
	for !win.Closed() {
		win.Clear(colornames.Black)
		log_fps(&last_time, particles)
		key_listener(win, &particles)
		if !PAUSED {
			if !RAN {
				kernel_call(context, queue, kernel)
				RAN = true
			}
			// update_particles(particles)
		}
		draw_particles(win, particles)
		win.Update()
	}
}

func log_fps(last_time *time.Time, particles []Particle) {
	current_time := time.Now()
	elapsed := current_time.Sub(*last_time).Seconds()
	*last_time = current_time
	fps := 1.0 / elapsed
	fmt.Printf("\rParticle Count: %d | FPS: %.2f", len(particles), fps)
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
