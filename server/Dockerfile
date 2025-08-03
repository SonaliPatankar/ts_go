# ---------- Stage 1: Build ----------
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything else
COPY . .

# ⚠️ Build statically (no glibc dependency)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main


# ---------- Stage 2: Lambda Runtime ----------
FROM public.ecr.aws/lambda/go:1

# Copy the binary from builder
COPY --from=builder /app/main .

# Lambda entrypoint
CMD ["main"]
