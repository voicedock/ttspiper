FROM golang:1.19-alpine as golang
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

FROM bufbuild/buf
COPY --from=golang /go/bin/protoc-gen-go-grpc /go/bin/protoc-gen-go-grpc
COPY --from=golang /go/bin/protoc-gen-go /go/bin/protoc-gen-go
ENV PATH="/go/bin:${PATH}"
ENTRYPOINT ["/usr/local/bin/buf"]