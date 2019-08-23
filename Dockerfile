FROM golang as builder

WORKDIR /go/src/cloudrun/line
COPY . .
RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -v -o line

FROM marketplace.gcr.io/google/ubuntu1804:latest

COPY --from=builder /go/src/cloudrun/line/line /line

CMD ["/line"]
