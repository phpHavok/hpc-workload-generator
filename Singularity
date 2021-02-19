Bootstrap: library
From: ubuntu:20.04

%environment
    export PATH=$PATH:/usr/local/go/bin

%post
    apt-get -y update
    apt-get -y install htop wget build-essential
    TMPGO=`mktemp`
    wget -c https://golang.org/dl/go1.16.linux-amd64.tar.gz -O $TMPGO
    tar -C /usr/local -xzf $TMPGO
    rm -f $TMPGO
