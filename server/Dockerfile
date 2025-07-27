# Stage 1: Build the Go binary
FROM golang:1.21 AS build

WORKDIR /app

# Copy go.mod and go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build it
COPY . .
RUN go build -o main

# Stage 2: Create the Lambda image
FROM public.ecr.aws/lambda/go:1

# Copy compiled binary from build stage
COPY --from=build /app/main .

# Set the CMD to run the binary
CMD ["main"]
