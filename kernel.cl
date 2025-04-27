__kernel void update_velocities_and_positions(__global float* positions_and_velocities,
                                              __global float* output,
                                              __global const float* interactionMatrix,
                                              const int particle_count, 
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
    
    float vx = positions_and_velocities[5 * i + 2];
    float vy = positions_and_velocities[5 * i + 3];

    for (int j = i; j < particle_count; ++j) {
        if (i == j) continue;

        float tx = positions_and_velocities[5 * j + 2];
        float ty = positions_and_velocities[5 * j + 3];

        // Particle B position
        float qx = positions_and_velocities[5 * j];
        float qy = positions_and_velocities[5 * j + 1];

        // Calculate distance and force (acceleration)
        float dx = px - qx;
        float dy = py - qy;
        float r2 = dx * dx + dy * dy;
        float r = sqrt(r2);

        float inverse = 0.0f;
        float damp = 1.0f;

        if (r > FORCE_RANGE) {
            inverse = 0.0f;
            damp = 1.0f;
        } else if (r < 1.0f) {
            inverse = -1.0f / (r + EPSILON);
            damp = PROXIMAL_DAMP;
        } else {
            inverse = 1.0f / r;
            damp = DISTAL_DAMP;
        }

        float ax = CHARGE * dx * inverse;
        float ay = CHARGE * dy * inverse;

        float if_a = interactionMatrix[i * 5 + 4];
        float if_b = interactionMatrix[j * 5 + 4];

        // Update velocity for Particle A
        vx = (vx + ax * if_a) * damp;
        vy = (vy + ay * if_a) * damp;

        tx = (tx + -ax * if_b) * damp;
        ty = (ty + -ay* if_b) * damp;

        output[5 * i + 2] = vx;
        output[5 * i + 3] = vy;
        output[5 * j + 2] = tx;
        output[5 * j + 3] = ty;
    }

    // Update positions based on calculated velocities
    float new_x = px + vx * DELTA;
    float new_y = py + vy * DELTA;

    // Wall collision detection
    if (new_x < X_MIN_BOUND) {
        new_x = X_MIN_BOUND;
    } else if (new_x > X_MAX_BOUND) {
        new_x = X_MAX_BOUND;
    }

    if (new_y < Y_MIN_BOUND) {
        new_y = Y_MIN_BOUND;
    } else if (new_y > Y_MAX_BOUND) {
        new_y = Y_MAX_BOUND;
    }

    // Store the new position back to the positions array
    output[5 * i] = new_x;
    output[5 * i + 1] = new_y;

    if (output[5 * i] < X_MIN_BOUND) output[5 * i] = X_MIN_BOUND;
    if (output[5 * i] > X_MAX_BOUND) output[5 * i] = X_MAX_BOUND;
    if (output[5 * i + 1] < Y_MIN_BOUND) output[5 * i + 1] = Y_MIN_BOUND;
    if (output[5 * i + 1] > Y_MAX_BOUND) output[5 * i + 1] = Y_MAX_BOUND;
}
