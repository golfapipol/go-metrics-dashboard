FROM golang:1.12 as builder
WORKDIR /module
COPY ./go.mod /module/go.mod
COPY ./go.sum /module/go.sum
RUN go mod download
COPY ./cmd /module/cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/main.go

FROM alpine:3.9
WORKDIR /root/
COPY --from=builder /module/bin .
ENV GIN_MODE release
EXPOSE 8000
CMD ./app