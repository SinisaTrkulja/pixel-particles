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
                                              }