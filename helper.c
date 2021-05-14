#define _GNU_SOURCE
#include <stdlib.h>
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

// Allocate memory of size bytes and initialize it so that the OS marks it as
// actually used
void *allocate_memory(unsigned long bytes)
{
    unsigned long i = 0;
    void *ptr = malloc(bytes);
    if (NULL == ptr)
    {
        return NULL;
    }
    // For some reason, memset doesn't actually seem to mark the memory as
    // "in-use", so we manually fill it with zeros
    for (i = 0; i < bytes; i++)
    {
        *((char *)ptr + i) = 0;
    }
    return ptr;
}

void release_memory(void *ptr)
{
    free(ptr);
}