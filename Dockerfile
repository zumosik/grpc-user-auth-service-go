# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

COPY app ./
RUN go build -o ./bin/app cmd/main.go

FROM alpine AS runner

# Expose the gRPC server port
EXPOSE 8081

COPY --from=builder /usr/local/src/bin/app /
COPY configs/cfg.yml /configs/cfg.yml

CMD ["/app"]
