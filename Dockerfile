# Stage 1: Build frontend assets
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build the Go binary with embedded static assets
FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist

# Build the binary statically so it runs in distroless/scratch
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seraphine .

# Stage 3: Distroless runtime image
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/seraphine .

EXPOSE 9009
EXPOSE 8080
ENTRYPOINT ["/app/seraphine", "server"]
