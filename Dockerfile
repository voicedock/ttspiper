FROM golang:1.20 as builder

ARG TARGETARCH
ARG TARGETVARIANT

RUN apt update && apt install -y \
        g++ \
        automake \
        cmake \
        pkg-config \
        libtool \
        gcc libc-dev binutils-gold \
        wget \
        git && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

ARG SPDLOG_VERSION="1.11.0"
RUN curl -L "https://github.com/gabime/spdlog/archive/refs/tags/v${SPDLOG_VERSION}.tar.gz" | \
    tar -xzvf - && \
    mkdir -p "spdlog-${SPDLOG_VERSION}/build" && \
    cd "spdlog-${SPDLOG_VERSION}/build" && \
    cmake -DCMAKE_POSITION_INDEPENDENT_CODE=ON ..  && \
    make -j8 && \
    cmake --install . --prefix /usr

# Use pre-compiled Piper phonemization library (includes onnxruntime)
ARG PIPER_PHONEMIZE_VERSION='1.1.0'
RUN mkdir -p "lib/Linux-$(uname -m)/piper_phonemize" && \
    curl -L "https://github.com/rhasspy/piper-phonemize/releases/download/v${PIPER_PHONEMIZE_VERSION}/libpiper_phonemize-${TARGETARCH}${TARGETVARIANT}.tar.gz" | \
        tar -C "lib/Linux-$(uname -m)/piper_phonemize" -xzvf -

ADD cpp_src/include /build/include
ADD cpp_src/src /build/src
ADD cpp_src/CMakeLists.txt /build/CMakeLists.txt
ADD . /usr/src/app

RUN mkdir -p build && \
    cd build && \
    cmake .. -DCMAKE_BUILD_TYPE=Release && \
    make install && \
    cd .. && \
    mkdir /usr/src/app/lib && \
    export PP_DIR="lib/Linux-$(uname -m)/piper_phonemize" && \
    cp -aR ${PP_DIR}/lib/espeak-ng-data ${PP_DIR}/lib/*.so* /usr/local/lib/libttspiperlib* /usr/src/app/lib/

RUN cd /usr/src/app && \
    cp ./cpp_src/include/ttspiperlib/ttspiperlib.h /usr/local/include && \
    export CGO_LDFLAGS="-L/usr/src/app/lib -Wl,-rpath -Wl,\$ORIGIN/../lib" && \
    export CGO_CFLAGS="-I/usr/local/include" && \
    rm -rf ~/.cache/go-build && \
    go build -o ./ttspiper ./cmd/ttspiper && \
    tar -czf lib.tar.gz lib/

FROM debian:12

WORKDIR /usr/src/app

RUN apt update && \
    apt install -y ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/lib.tar.gz ./
COPY --from=builder /usr/src/app/ttspiper /usr/src/app/ttspiper

RUN tar -xzf lib.tar.gz && mv ./lib/espeak-ng-data /usr/share/espeak-ng-data

ENV PATH=${PATH}:/usr/src/app/

CMD ["ttspiper"]