FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go build -o server

FROM alpine:latest
RUN apk add --no-cache ffmpeg python3 curl \
    && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
