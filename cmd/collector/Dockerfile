FROM golang:1.18 as builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/app .

FROM alpine:latest
COPY --from=builder /usr/local/bin/app ./
CMD ["./app"]