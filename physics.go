package main

import (
	"math"
	"sync"
)

func update_velocities(particles []Particle) {
	var wg sync.WaitGroup
	for i := range len(particles) - 2 {
		pa := &particles[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := i + 1; j < len(particles); j++ {
				pb := &particles[j]
				update_velocity(pa, pb)
			}
			detect_wall_collision(pa)
		}()
	}
	wg.Wait()
	detect_wall_collision(&particles[len(particles)-1])
}

func update_velocity(pa, pb *Particle) {
	if_a, if_b := interaction(*pa, *pb)
	ax, ay, damp := calculate_acceleration(*pa, *pb)
	pa.x_speed = (pa.x_speed + ax*if_a) * damp
	pa.y_speed = (pa.y_speed + ay*if_a) * damp

	pb.x_speed = (pb.x_speed + (-1)*ax*if_b) * damp
	pb.y_speed = (pb.y_speed + (-1)*ay*if_b) * damp
}

func detect_wall_collision(p *Particle) {
	if p.x_position <= X_MIN_BOUND {
		p.x_speed = math.Abs(p.x_speed)
	} else if p.x_position >= X_MAX_BOUND {
		p.x_speed = -math.Abs(p.x_speed)
	}
	if p.y_position <= Y_MIN_BOUND {
		p.y_speed = math.Abs(p.x_speed)
	} else if p.y_position >= Y_MAX_BOUND {
		p.y_speed = -math.Abs(p.y_speed)
	}
}

func calculate_acceleration(acted_upon, acting Particle) (float64, float64, float64) {
	dx := acted_upon.x_position - acting.x_position
	dy := acted_upon.y_position - acting.y_position
	dx_2, dy_2 := dx*dx, dy*dy
	r := math.Sqrt(dx_2 + dy_2)
	inverse, damp := get_attraction(r)
	return CHARGE * dx * inverse, CHARGE * dy * inverse, damp
}

func get_attraction(r float64) (float64, float64) {
	if r > FORCE_RANGE {
		return 0.0, 1.0
	}
	if r < 1.0 {
		return -1.0 / (r + EPSILON), PROXIMAL_DAMP
	}
	return 1.0 / r, DISTAL_DAMP
}

func update_positions(particles []Particle) {
	var wg sync.WaitGroup
	for i := range particles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			update_position(&particles[i])
		}()
	}
	wg.Wait()
}

func update_position(particle *Particle) {
	new_x := particle.x_position + particle.x_speed*DELTA
	new_y := particle.y_position + particle.y_speed*DELTA
	if new_x < X_MIN_BOUND {
		particle.x_position = X_MIN_BOUND
	} else {
		particle.x_position = min(new_x, X_MAX_BOUND)
	}
	if new_y < Y_MIN_BOUND {
		particle.y_position = Y_MIN_BOUND
	} else {
		particle.y_position = min(new_y, Y_MAX_BOUND)
	}
}
