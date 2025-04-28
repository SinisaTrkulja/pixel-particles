__kernel void update_positions_and_velocities(__global const float* positions_and_velocities,
                                              __global float* output,
                                              __global const float* interactionMatrix,
                                              const float particle_count, 
                                              const float CHARGE, 
                                              const float FORCE_RANGE, 
                                              const float EPSILON, 
                                              const float PROXIMAL_DAMP, 
                                              const float DISTAL_DAMP,
                                              const float DELTA, 
                                              const float X_MIN_BOUND, const float X_MAX_BOUND,
                                              const float Y_MIN_BOUND, const float Y_MAX_BOUND) {
    int i = get_global_id(0);

    if (i >= particle_count) return;

    // Particle A position
    float px = positions_and_velocities[5 * i];
    float py = positions_and_velocities[5 * i + 1];

    // Accumulate acceleration for particle i
    float ax_total = 0.0f;
    float ay_total = 0.0f;
    int pos_a = positions_and_velocities[i*5 + 4];

    for (int j = 0; j < particle_count; j++) {
        if (i == j) continue;

        float qx = positions_and_velocities[5 * j];
        float qy = positions_and_velocities[5 * j + 1];

        int pos_b = positions_and_velocities[j*5 + 4];
        float if_a = interactionMatrix[pos_b*5 + pos_a]; // Use correct indexing

        float dx = px - qx;
        float dy = py - qy;
        float r2 = dx * dx + dy * dy;
        float r = sqrt(r2); 

        float ax = CHARGE * dx * ((r < FORCE_RANGE) ? ((r < 1.0f) ? -1.0f / (r  + EPSILON) : 1.0f / r) : 0.0f);
        float ay = CHARGE * dy * ((r < FORCE_RANGE) ? ((r < 1.0f) ? -1.0f / (r  + EPSILON) : 1.0f / r) : 0.0f);


        ax_total += ax * if_a * ((r < 1.0f) ? PROXIMAL_DAMP : DISTAL_DAMP);
        ay_total += ay * if_a * ((r < 1.0f) ? PROXIMAL_DAMP : DISTAL_DAMP);
    }

    // Update velocity
    float vx = positions_and_velocities[5 * i + 2] + ax_total;
    float vy = positions_and_velocities[5 * i + 3] + ay_total;

    

    // Update positions based on calculated velocities
    float new_x = px + vx * DELTA;
    float new_y = py + vy * DELTA;

    // Wall collision detection
    if (new_x < X_MIN_BOUND || new_x > X_MAX_BOUND) {
        vx *= -1;
    } 
    if (new_y < Y_MIN_BOUND || new_y > Y_MAX_BOUND) {
        vy *= -1;
    } 

    // Store the new velocities
    output[5 * i + 2] = vx;
    output[5 * i + 3] = vy;

    // Store the new positions
    if (new_x < X_MIN_BOUND) {
        output[5 * i] = X_MIN_BOUND;
    } else if (new_x > X_MAX_BOUND) {
        output[5 * i] = X_MAX_BOUND;
    } else {
        output[5 * i] = new_x;
    }

    if (new_y < Y_MIN_BOUND) {
        output[5 * i + 1] = Y_MIN_BOUND;
    } else if (new_y > Y_MAX_BOUND) {
        output[5 * i + 1] = Y_MAX_BOUND;
    } else {
        output[5 * i + 1] = new_y;
    }
}
