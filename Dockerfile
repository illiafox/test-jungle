# Initial stage: download modules
FROM golang:1.20 as modules

WORKDIR /app
ADD app/go.mod go.mod
ADD app/go.sum go.sum

RUN --mount=type=cache,target=/go/pkg \
     cd /app && go mod download

# Intermediate stage: Build the binary
FROM golang:1.20 as builder
COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /app
COPY ./app /app
WORKDIR /app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags '-s -w -extldflags "-static"' \
    -o ./bin/main ./cmd/server/main.go

# Final stage: Run the binary
FROM alpine:latest as image
COPY --from=builder /app/bin/main /main

EXPOSE 8080
EXPOSE 8082
ENTRYPOINT [ "/main" ]