FROM golang:1.20 as builder

RUN apt update && apt install -y \
        g++ \
        automake \
        cmake \
        pkg-config \
        libtool \
        wget \
        git && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

ADD cpp_src/lib /build/lib
RUN tar -xvf "lib/espeak-ng-1.52-patched.tar.gz" -C "./" && \
    cd espeak-ng && \
    ./autogen.sh && \
    ./configure \
        --without-pcaudiolib \
        --without-klatt \
        --without-speechplayer \
        --without-mbrola \
        --without-sonic \
        --with-extdict-cmn \
        --prefix=/usr && \
    make -j8 src/espeak-ng src/speak-ng && \
    make && \
    make install

RUN export ONNX_DIR="./lib/Linux-$(uname -m)" && \
    export ONNX_FILENAME="onnxruntime-linux-x64-1.14.1.tgz" && \
    wget "https://github.com/microsoft/onnxruntime/releases/download/v1.14.1/${ONNX_FILENAME}" && \
    mkdir -p "${ONNX_DIR}" && \
    tar -C "${ONNX_DIR}" \
        --strip-components 1 \
        -xvf "${ONNX_FILENAME}"

ADD cpp_src/include /build/include
ADD cpp_src/src /build/src
ADD cpp_src/CMakeLists.txt /build/CMakeLists.txt

RUN mkdir build && cd build/ && cmake .. && make install

FROM golang:1.20

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

WORKDIR /usr/src/app
COPY --from=builder /usr/local/lib/libttspiperlib.so /usr/local/lib/
COPY --from=builder /usr/local/include/ttspiperlib.h /usr/local/include
COPY --from=builder /usr/lib/libespea* /usr/lib/
COPY --from=builder /usr/share/espeak-ng-data /usr/share/espeak-ng-data
COPY --from=builder /build/lib/Linux-x86_64/lib/libonnxruntime* /usr/local/lib/

ADD . /usr/src/app

RUN go build -o ./ttspiper ./cmd/ttspiper && ldconfig

CMD ["sttwhisper"]