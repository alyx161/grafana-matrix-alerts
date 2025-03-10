##
# Golang Builder
##
FROM golang:latest AS builder

# Install necessary dependencies
RUN apt update && apt install -y git libolm-dev

# Copy src
COPY . /go/src/gma

# Set workdir
WORKDIR /go/src/gma

# Build
RUN go generate
RUN go mod tidy
RUN go build -o gma


##
# Main image
##
FROM debian:stable-slim
RUN apt update && apt upgrade -y && apt install -y git libolm-dev

# Copy build
COPY --from=builder /go/src/gma/gma /app/gma
RUN chmod +x /app/gma

# Set workdir
WORKDIR /app

# Expose ports
EXPOSE 8080

CMD ["/app/gma"]