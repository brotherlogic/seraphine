FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build the binary statically so it runs in distroless/scratch
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seraphine .

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/seraphine .

EXPOSE 9009
ENTRYPOINT ["/app/seraphine", "server"]
