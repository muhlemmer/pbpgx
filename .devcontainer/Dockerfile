FROM golang:1.19

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
        git ssh-client autoconf automake libtool curl make g++ unzip \
    && rm -rf /var/lib/apt/lists/*

ARG PROTOBUF_VERSION=3.19.1 JOBS=8

RUN git clone --branch v${PROTOBUF_VERSION} --depth 1 --recurse-submodules --shallow-submodules https://github.com/protocolbuffers/protobuf.git /protobuf && \
    cd /protobuf && \
    ./autogen.sh && \
    ./configure && \
    make -j${JOBS} && \
   # make -j${JOBS} check && \
    make install && \
    ldconfig

ARG USER=vscode
ARG UID=1000
ARG GROUP=vscode
ARG GID=1000

RUN groupadd -g ${GID} ${GROUP}
RUN useradd -m -g ${GID} -s /bin/bash ${USER}

USER ${USER}
RUN go install github.com/cweill/gotests/gotests@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/mgechev/revive@latest && \
    go install github.com/ramya-rao-a/go-outline@latest && \
    go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install golang.org/x/lint/golint@latest && \
    go install golang.org/x/tools/cmd/godoc@master && \
    go install golang.org/x/tools/gopls@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest

# godoc server
EXPOSE 6060
