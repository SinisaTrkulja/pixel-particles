package main

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
