# official Go image
FROM golang:1.24-alpine AS builder

# set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first. This is a Docker caching trick.
# If these files don't change, Docker will use the cached dependencies layer,
# speeding up subsequent builds.
# go.sum is the lock file for dependency management, 
# this project only uses standard lib so there is no sum file so the * skips it
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# build 
# -o /server: specifies the output file name and location.
# -ldflags="-w -s": makes the binary smaller by stripping debug information
# CGO_ENABLED=0: creates a statically linked binary (important for minimal container)
RUN CGO_ENABLED=0 go build -o /server -ldflags="-w -s" ./cmd/server/main.go


#final image, slim base
FROM alpine:latest

# only need to copy the compiled binary from the builder stage
# --from=builder flag
COPY --from=builder /server /server

# copy web templates 
COPY ./web/template /web/template

# EXPOSE tells Docker which port our application listens on.
# This is mainly for documentation; it doesn't actually open the port.
EXPOSE 8080

# CMD defines the command that will run when the container starts.
CMD ["/server"]