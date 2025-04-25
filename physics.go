package main

import (
	"math"
	"sync"
)

func calcute_velocities(particles []Particle) {
	var wg sync.WaitGroup
	for i := range particles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acted_upon := &particles[i]
			for j := range particles {
				if i != j {
					acting := particles[j]
					sign := interaction_sign(*acted_upon, acting)
					ax, ay, damp := calculate_acceleration(*acted_upon, acting)
					acted_upon.x_speed = (acted_upon.x_speed + ax*sign) * damp
					acted_upon.y_speed = (acted_upon.y_speed + ay*sign) * damp
				}
			}
		}()
	}
	wg.Wait()
}

func calculate_acceleration(acted_upon, acting Particle) (float64, float64, float64) {
	dx := acted_upon.x_position - acting.x_position
	dy := acted_upon.y_position - acting.y_position

	dx_2 := dx * dx
	dy_2 := dy * dy
	r := math.Sqrt(dx_2 + dy_2)
	inverse, damp := get_attraction(r)
	return CHARGE * dx * inverse, CHARGE * dy * inverse, damp
}

func get_attraction(r float64) (float64, float64) {
	if r > FORCE_RANGE {
		return 0.0, 1.0
	} else if r < 1.0 {
		return -1.0 / (r + EPSILON), PROXIMA_DAMP
	} else {
		return 1.0 / r, DISTAL_DAMP
	}
}

func update_positions(particle *Particle) {
	x := particle.x_position
	y := particle.y_position

	if x <= X_MIN+RADIUS || x >= X_MAX-RADIUS {
		particle.x_speed *= -1
	}
	if y <= Y_MIN+RADIUS || y >= Y_MAX-RADIUS {
		particle.y_speed *= -1
	}

	particle.x_position += particle.x_speed * DELTA
	particle.y_position += particle.y_speed * DELTA
}
