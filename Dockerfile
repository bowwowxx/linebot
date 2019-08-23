FROM golang as builder

WORKDIR /go/src/cloudrun/line
COPY . .
RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -v -o line

FROM marketplace.gcr.io/google/ubuntu1804:latest

ENV HostPort :8080
ENV LineSecret c003487f9082e4c1aa9cfcea177c8d51
ENV LineToken wuHvQ9og6UDT6CAUHrLy7Fuf5ekznyVcAOpuFds2X3Vnyp0Mf65BJvGPCpOb/drO+HfFTy6aSSr3xQGLbeOZfMNXr6hdAxMzcSWx4hO9QbSkIYVteowk0JfP3rBrO48rUQQK5qfPLAbfpk0sea5EmAdB04t89/1O/w1cDnyilFU=

COPY --from=builder /go/src/cloudrun/line/line /line

CMD ["/line"]
