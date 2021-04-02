# HPC Workload Generator

A Go project for generating defined HPC workload characteristics (CPU usage, RAM
usage, etc.) on a system to facilitate granular monitoring.

## Compiling from Source

This project is written primarily in Go and requires Go v1.15 or later to
compile. Cgo is also used, so you'll need a recent version of gcc.

To build, you can just type `go build`, and Go will handle everything.
Alternatively, if you have Make installed, you can just type `make`. Both of
these methods produce an executable binary called `hpc-workload-generator`.

## Usage

To view help, run `./hpc-workload-generator -h`.

```
$ ./hpc-workload-generator -h
Usage of ./hpc-workload-generator:
  -cgroups-root string
        path to the root of the cgroupsv1 hierarchy (default "/sys/fs/cgroup")
  -d uint
        debug level (0-6) where higher numbers have higher verbosity (default 4)
  -i string
        path to the schedule file to execute, or empty for stdin
```

The `-i` option is the most critical and specifies the path to the input file
from which to read the schedule of events. If not specified, the program will
read from standard input.

The `-cgroups-root` option allows you to change the default location of the
cgroupsv1 hierarchy if it happens to be mounted somewhere unusual. This can also
be handy when using a container in case you want to mount the hierarchy in a
different location for the sake of the container.

The `-d` option allows you to set the debug level. The default (4) is fine for
most cases, but you can bump it up to 5 or 6 if you want more verbose logging.

## Docker

A Docker container is provided for systems where that is more convenient (such
as a Kubernetes cluster). You can build it manually using the provided
`Dockerfile`, or just pull the pre-built copy from Docker Hub. Example usage
follows:

```
docker run -t --rm \
    --mount type=bind,src=/sys/fs/cgroup,dst=/sys/fs/cgroup,readonly \
    --mount type=bind,src=`pwd`,dst=/data \
    phphavok/hpc-workload-generator -i /data/example.schedule
```

We specify `-t` so that we're allocated a pseudo-terminal which makes the
logging output look nice and formatted. The `--rm` option automatically cleans
up the container on exit. The first mount command passes through the cgroupv1
hierarchy (/sys/fs/cgroup) on the parent system to the same location within the
container. By default, Docker will often have some of the cgroup hierarchy
present within the container, but not all of it. This application will need to
see the full hierarchy, so this read-only bind mount takes care of that. If you
run into issues with mounting over the existing hierarchy within the container,
you can change the target to some other location and then pass the
`-cgroup-root` option to the program to accommodate that change. The second
mount command just passes through the current working directory (assumed to be
this repository's root) so that schedule files can be accessed and ran on from
within the container. The entrypoint to the container is the program itself, so
you can just pass its parameters to the run command.

## Singularity

You can use Singularity (e.g., on an HPC cluster) to run the above Docker
container from Docker Hub.

```
singularity run \
    -B /sys/fs/cgroup:/sys/fs/cgroup \
    -B `pwd`:/data \
    docker://phphavok/hpc-workload-generator -i /data/example.schedule
```

## Schedule of Events

The schedule of events is an input file that the hpc-workload-generator program
reads. The format of the file is one event per line, where an event line is a
comma-separated list of values. The first value indicates at how many seconds
into the runtime the event should fire. The second value is the name of the
event module to run. All remaining values are arguments to the specified event
module. An example schedule of events follows with commentary:

```
4, cpuload, 1, 70, 12
8, cpuload, 3, 70, 12
12, cpuload, 5, 70, 12
16, cpuload, 7, 70, 12
```

This schedule will execute the cpuload module four different times at 4 seconds,
8 seconds, 12 seconds, and 16 seconds, respectively. Each execution has
different potential settings. Note that modules can and will be executed in
parallel if you specify multiple modules at the same second offset (or if one
module is still running when a later module begins).

This particular schedule will cause CPU number 1 (0-indexed) to be taxed at
roughly 70% for 12 seconds starting at 4 seconds into the runtime of the
program. At 8 seconds into the runtime, CPU number 3 will start being taxed at
70% for 12 seconds. Note that at this point, only 4 seconds have passed since
CPU 1 started being taxed, so it still has 8 seconds left of being taxed. Thus,
CPU 1 and 3 will be taxed in parallel at this point. The rest of the schedule
follows a similar scheme.

Read more about the supported modules and their parameters below.

### cpuload

The `cpuload` module generates simulated load on a CPU. The module has three
arguments.

The first argument is the index of the CPU to generate load on. It is 0-indexed
while programs like `htop` display the CPUs as 1-indexed, so you'll have to do
some mental gymnastics when looking at live results. Furthermore, the CPU index
is automatically mapped into the program's active cgroup. Suppose you're running
the program on a node that has 16 CPUs, 8 of which have been reserved for your
node via a cgroup. It's possible that you'll be reserved CPUs that aren't
contiguous. That is, instead of getting CPUs 0 through 7, you may get CPUs 0-3
and then 8 through 11. You won't know this ahead of time when writing a schedule
of events though. Therefore, just specify CPU indices as though they've been
reserved contiguously (i.e., 0 through 7), and if this is not the case, the
program will automatically map your specified CPU indices to the real indices
that were reserved.

The second argument is the decimal percentage (1-100) of load you'd like to
generate on that CPU. It will generate load close to this number. However, if
you generate multiple loads on the same CPU index in parallel, it may not do
what you'd expect as each of the load generators don't know about each other and
will independently assume they have total control over that CPU.

The third and final argument is the number of seconds for which the load should
be generated.
