FROM golang:1.20-alpine
LABEL org.opencontainers.image.source=https://github.com/barpav/msg-files
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN ping -w 1 proxy.golang.org \
    && go mod download && go mod verify
COPY cmd/ cmd/
COPY internal/ internal/
COPY statistics/ statistics/
RUN go build -v -o /usr/local/bin/app ./cmd
CMD ["app"]