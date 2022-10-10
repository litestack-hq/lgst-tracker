FROM golang:1.18.0-buster as builder
RUN go env -w GOPRIVATE=github.com/litestack-hq/
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN --mount=type=ssh go mod download
COPY . .
RUN go generate ./...
RUN go build -o .out/app
FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/.out .out
CMD ./.out/app
