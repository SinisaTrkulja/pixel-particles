package main

import (
	"fmt"
	"log"
	"sync"
	"unsafe"

	"github.com/jgillich/go-opencl/cl"
)

const kernel_source = "kernel.cl"

func init_kernel() (*cl.Context, *cl.CommandQueue, *cl.Kernel) {

	// Initialize OpenCL
	platform, err := cl.GetPlatforms()
	if err != nil {
		log.Fatal(err)
	}

	device, err := platform[0].GetDevices(cl.DeviceTypeGPU)
	if err != nil {
		log.Fatal(err)
	}

	context, err := cl.CreateContext([]*cl.Device{device[0]})
	if err != nil {
		log.Fatal(err)
	}

	queue, err := context.CreateCommandQueue(device[0], 0)
	if err != nil {
		log.Fatal(err)
	}

	program, err := context.CreateProgramWithSource([]string{kernel_source})
	if err != nil {
		log.Fatal(err)
	}

	err = program.BuildProgram(nil, "")
	if err != nil {
		log.Fatal(err)
	}

	kernel, err := program.CreateKernel("update_velocities_and_positions")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Kernel initialized successfully")
	return context, queue, kernel
}

func kernel_call(context *cl.Context, queue *cl.CommandQueue, kernel *cl.Kernel) {

	posBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, int(unsafe.Sizeof(float32(0)))*len(positions_and_velocities))
	if err != nil {
		log.Fatal(err)
	}
	defer posBuffer.Release()

	fmt.Println("Position buffer created successfully")

	// Create buffer for interaction matrix and particle colors
	interactionMatrix2Buffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, int(unsafe.Sizeof(float32(0)))*len(interactionMatrix2))
	if err != nil {
		log.Fatal(err)
	}
	defer interactionMatrix2Buffer.Release()

	fmt.Println("interaction buffer created successfully")

	// Write data to buffers
	positionElemSize := int(unsafe.Sizeof(positions_and_velocities[0]))
	inputDataPtr := unsafe.Pointer(&positions_and_velocities[0])
	inputDataTotalSizeBytes := positionElemSize * len(positions_and_velocities)
	_, err = queue.EnqueueWriteBuffer(posBuffer, true, 0, inputDataTotalSizeBytes, inputDataPtr, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Position write buffer enqueued successfully")

	ifElemSize := int(unsafe.Sizeof(interactionMatrix2[0]))
	ifDataPtr := unsafe.Pointer(&interactionMatrix2[0])
	ifDataTotalSizeBytes := ifElemSize * len(interactionMatrix2)
	_, err = queue.EnqueueWriteBuffer(interactionMatrix2Buffer, true, 0, ifDataTotalSizeBytes, ifDataPtr, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matrix write buffer enqueued successfully")

	outputBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, positionElemSize*len(positions_and_velocities))
	if err != nil {
		log.Fatal(err)
	}
	defer outputBuffer.Release()

	fmt.Println("CreateEmptyBuffer successfully")

	// Set kernel arguments
	err = kernel.SetArgs(posBuffer, outputBuffer, interactionMatrix2Buffer, float32(PARTICLE_COUNT),
		float32(CHARGE), float32(FORCE_RANGE), float32(EPSILON), float32(PROXIMAL_DAMP), float32(DISTAL_DAMP), float32(DELTA), float32(X_MIN_BOUND),
		float32(X_MAX_BOUND), float32(Y_MIN_BOUND), float32(Y_MAX_BOUND))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SetArgs successfully")

	// Execute kernel
	globalWorkSize := []int{PARTICLE_COUNT}
	_, err = queue.EnqueueNDRangeKernel(kernel, nil, globalWorkSize, []int{8}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("EnqueueNDRangeKernel successfully")

	err = queue.Finish()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished successfully")

	results := make([]float32, len(positions_and_velocities))
	dataPtrOut := unsafe.Pointer(&results[0])
	sizePerEntry := int(unsafe.Sizeof(results[0]))
	dataSizeOut := sizePerEntry * len(results)

	_, err = queue.EnqueueReadBuffer(outputBuffer, true, 0, dataSizeOut, dataPtrOut, nil)

	fmt.Println("EnqueueReadBuffer successfully")

	if err != nil {
		log.Fatal(err)
	}

	// Output results
	fmt.Println("Update positions:", positions_and_velocities)
	wg := sync.WaitGroup{}
	for i := range PARTICLE_COUNT {
		wg.Add(1)
		go func(i int) {
			particles[i].x_position = float64(results[i*5])
			particles[i].y_position = float64(results[i*5+1])
			particles[i].x_speed = float64(results[i*5+2])
			particles[i].y_speed = float64(results[i*5+3])
		}(i)
	}
	wg.Wait()
}
