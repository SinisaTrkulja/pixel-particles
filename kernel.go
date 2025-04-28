package main

import (
	"fmt"
	"log"
	"sync"
	"unsafe"

	"github.com/jgillich/go-opencl/cl"
)

const kernel_source = "kernel2.cl"

func kernel_call(particles []Particle) {

	// Initialize OpenCL
	platform, err := cl.GetPlatforms()
	if err != nil {
		log.Fatal(err)
	}

	device, err := platform[0].GetDevices(cl.DeviceTypeGPU)
	if err != nil {
		log.Fatal(err)
	}

	for i, d := range device {
		fmt.Println(i, d.Name())
	}

	var device_id int
	if len(device) >= 1 {
		fmt.Print("Choose your device: ")
		_, err = fmt.Scan(&device_id)
		if err != nil {
			log.Fatal("Failed to read input:", err)
		}

		if device_id < 0 || device_id >= len(device) {
			log.Fatal("Invalid device ID")
		}
	}

	context, err := cl.CreateContext([]*cl.Device{device[device_id]})
	if err != nil {
		log.Fatal(err)
	}

	queue, err := context.CreateCommandQueue(device[0], 0)
	if err != nil {
		log.Fatal(err)
	}
	defer queue.Release()

	program, err := context.CreateProgramWithSource([]string{kernel_source})
	if err != nil {
		log.Fatal(err)
	}

	err = program.BuildProgram(nil, "")
	if err != nil {
		log.Fatal(err)
	}

	velocities_kernel, err := program.CreateKernel("update_positions_and_velocities")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Kernel initialized successfully")

	posBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, int(unsafe.Sizeof(float32(0)))*len(positions_and_velocities))
	if err != nil {
		log.Fatal(err)
	}
	defer posBuffer.Release()

	// fmt.Println("Position buffer created successfully")

	// Create buffer for interaction matrix and particle colors
	interactionMatrix2Buffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, int(unsafe.Sizeof(float32(0)))*len(interactionMatrix2))
	if err != nil {
		log.Fatal(err)
	}
	defer interactionMatrix2Buffer.Release()

	// fmt.Println("interaction buffer created successfully")

	positionElemSize := int(unsafe.Sizeof(positions_and_velocities[0]))
	outputBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, positionElemSize*len(positions_and_velocities))
	if err != nil {
		log.Fatal(err)
	}
	defer outputBuffer.Release()

	inputDataPtr := unsafe.Pointer(&positions_and_velocities[0])
	inputDataTotalSizeBytes := positionElemSize * len(positions_and_velocities)
	_, err = queue.EnqueueWriteBuffer(posBuffer, true, 0, inputDataTotalSizeBytes, inputDataPtr, nil)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Position write buffer enqueued successfully")

	ifElemSize := int(unsafe.Sizeof(interactionMatrix2[0]))
	ifDataPtr := unsafe.Pointer(&interactionMatrix2[0])
	ifDataTotalSizeBytes := ifElemSize * len(interactionMatrix2)
	_, err = queue.EnqueueWriteBuffer(interactionMatrix2Buffer, true, 0, ifDataTotalSizeBytes, ifDataPtr, nil)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Matrix write buffer enqueued successfully")

	// Set kernel arguments
	err = velocities_kernel.SetArgs(posBuffer, outputBuffer, interactionMatrix2Buffer,
		float32(PARTICLE_COUNT), float32(CHARGE), float32(FORCE_RANGE),
		float32(EPSILON), float32(PROXIMAL_DAMP), float32(DISTAL_DAMP),
		float32(DELTA), float32(X_MIN_BOUND), float32(X_MAX_BOUND),
		float32(Y_MIN_BOUND), float32(Y_MAX_BOUND))
	if err != nil {
		log.Fatal(err)
	}

	// Execute kernel
	globalWorkSize := []int{PARTICLE_COUNT}
	_, err = queue.EnqueueNDRangeKernel(velocities_kernel, nil, globalWorkSize, []int{COUNT_CHANGE_STEP}, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = queue.Finish()
	if err != nil {
		log.Fatal(err)
	}

	results := make([]float32, PARTICLE_COUNT*5) // to be fixed
	dataPtrOut := unsafe.Pointer(&results[0])
	sizePerEntry := int(unsafe.Sizeof(results[0]))
	dataSizeOut := sizePerEntry * len(results)

	_, err = queue.EnqueueReadBuffer(outputBuffer, true, 0, dataSizeOut, dataPtrOut, nil)
	if err != nil {
		log.Fatal(err)
	}

	update_particles_2(particles, results)

}

func update_particles_2(particles []Particle, results []float32) {
	wg := sync.WaitGroup{}
	for i := range PARTICLE_COUNT {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Update
			pa := &particles[i]
			pa.x_position = float64(results[i*5])
			pa.y_position = float64(results[i*5+1])
			pa.x_speed = float64(results[i*5+2])
			pa.y_speed = float64(results[i*5+3])

			// Sync
			positions_and_velocities[i*5] = float32(particles[i].x_position)
			positions_and_velocities[i*5+1] = float32(particles[i].y_position)
			positions_and_velocities[i*5+2] = float32(particles[i].x_speed)
			positions_and_velocities[i*5+3] = float32(particles[i].y_speed)
		}(i)
	}
	wg.Wait()
}
