#ifndef HELPER_H_
#define HELPER_H_

void lock_os_thread(int cpuid);
void *allocate_memory(unsigned long bytes);
void release_memory(void *ptr);

#endif