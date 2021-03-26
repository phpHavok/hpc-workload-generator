FROM golang:1.16-buster

RUN apt-get update && \
    apt-get install -y --no-install-recommends htop build-essential &&  \
    rm -rf /var/lib/apt/lists/*

COPY . /usr/local/hpc-workload-generator

RUN make -C /usr/local/hpc-workload-generator

ENTRYPOINT ["/usr/local/hpc-workload-generator/hpc-workload-generator"]
