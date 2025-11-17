FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod /src/go.mod
COPY hw3 /src/hw3
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /src/hw3/cmd/server && go build -o /out/balance-service

FROM ubuntu:22.04
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/balance-service /app/server
ENV LISTEN_ADDR=:8080
ENV STORAGE=memory
ENV FILE_PATH=/data/balances.json
EXPOSE 8080
ENTRYPOINT ["/app/server"]
