FROM golang:1.23 AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/server ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/server /app/server
ENV PORT=3000 APP_ENV=production
USER nonroot:nonroot
EXPOSE 3000
ENTRYPOINT ["/app/server"]
