package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
	"math/rand"
	"math"
	"image/color"
	"fmt"
)

const X_MIN = 0.0
const X_MAX = 1024.0
const Y_MIN = 0.0
const Y_MAX = 768.0

const SPEED = 4.0
const RADIUS = 4.0
const CHARGE = 20.0

const DELTA = .50

var POTENTIAL = 0.0


type Observable struct {
	x_position, y_position, x_speed, y_speed float64
	radius float64
	color pixel.RGBA                                          
}

func main() {
	opengl.Run(run)
}

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(X_MIN, Y_MIN, X_MAX, Y_MAX),
		VSync: true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	observables := init_observables(100)

	for !win.Closed() {
		win.Clear(colornames.Black)
		update_observables(win, observables)
		win.Update()
	}
}

func update_observables(win *opengl.Window, observables []Observable) {
	calcute_velocities(observables)
	draw_circles(win, observables)
	POTENTIAL = 0.0
}

func calcute_velocities(observables []Observable) {
	for i := 0; i < len(observables); i++ {
		acted_upon := &observables[i]
		for j := 0; j < len(observables); j++ {
			if i != j {
				acting := observables[j]
				ax, ay := calculate_acceleration(*acted_upon, acting)
				acted_upon.x_speed += ax // * DELTA
				acted_upon.y_speed += ay // * DELTA
			}
		}
	}
}

func calculate_acceleration(acted_upon, acting Observable) (float64, float64) {
	dx := acted_upon.x_position - acting.x_position
	dy := acted_upon.y_position - acting.y_position

	dx_2 := dx * dx
	dy_2 := dy * dy
	r := math.Sqrt(dx_2 + dy_2)
	inverse := 1.0/(r * r * r)

	POTENTIAL += CHARGE/r
	
	return CHARGE * dx * inverse, CHARGE * dy * inverse
}

func draw_circles(win *opengl.Window, observables []Observable) {
	total_kinetic := 0.0
	list_length := len(observables)
	for i := 0; i < list_length; i++ {
		observable := &observables[i]
		update_offsets(observable)

		imd := imdraw.New(nil)
		imd.Color = observable.color
		imd.Push(pixel.V(observable.x_position, observable.y_position))
		imd.Circle(observable.radius, 0)
		imd.Draw(win)

		x_speed := observable.x_speed
		y_speed := observable.y_speed
		total_kinetic += x_speed*x_speed + y_speed*y_speed
	}
	total_energy := total_kinetic + POTENTIAL
	fmt.Printf("\rtotal KE: %.3f, total PE: %.3f, total energy: %.3f, x-distance: %.3f", total_kinetic, POTENTIAL, total_energy, observables[0].x_position - observables[1].x_position)
}

func update_offsets(observable *Observable) {	
	x := observable.x_position
	y := observable.y_position
	radius := observable.radius

	if x - radius <= X_MIN || x + radius >= X_MAX {
		observable.x_speed *= -1
	}
	if y - radius <= Y_MIN || y + radius >= Y_MAX {
		observable.y_speed *= -1
	}

	observable.x_position += observable.x_speed * (DELTA / 2.0)
	observable.y_position += observable.y_speed * (DELTA / 2.0)
}

func init_observables(circle_count int) []Observable {
	var colors = []pixel.RGBA{
		to_pixel_color(colornames.Red),
		to_pixel_color(colornames.Blue),
		to_pixel_color(colornames.Green),
		to_pixel_color(colornames.Yellow),
		to_pixel_color(colornames.Purple),
	}
	observables := make([]Observable, circle_count)
	for i := 0; i < circle_count; i++ {
		observables[i] = Observable{
			x_position: rand.Float64() * float64(X_MAX),
			y_position: rand.Float64() * float64(Y_MAX),
			x_speed:  rand.Float64() * SPEED,
			y_speed:  rand.Float64() * SPEED,
			radius:   RADIUS,
			color: colors[rand.Intn(len(colors))],
		}
	}
	return observables
}


func to_pixel_color(c color.RGBA) pixel.RGBA {
	return pixel.RGBA{
		R: float64(c.R) / 255.0,
		G: float64(c.G) / 255.0,
		B: float64(c.B) / 255.0,
		A: float64(c.A) / 255.0,
	}
}