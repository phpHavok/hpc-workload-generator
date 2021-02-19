#define _GNU_SOURCE
#include <pthread.h>

// Affinity-binding code taken from: http://pythonwise.blogspot.com/2019/03/cpu-affinity-in-go.html
void lock_os_thread(int cpuid)
{
    pthread_t tid;
    cpu_set_t cpuset;
    tid = pthread_self();
    CPU_ZERO(&cpuset);
    CPU_SET(cpuid, &cpuset);
    pthread_setaffinity_np(tid, sizeof(cpu_set_t), &cpuset);
}