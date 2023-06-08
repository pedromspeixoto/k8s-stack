# syntax=docker/dockerfile:1

################################
# Build
################################
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN make build

################################
# Deploy
################################
FROM --platform=linux/amd64 debian:bullseye-slim

# Install Dependencies
RUN apt-get update  \
    && apt-get install -y ca-certificates wget xz-utils curl unzip python git \
    && update-ca-certificates

# Copy built binaries
RUN mkdir -p /opt/todo-api
COPY --from=build /app/bin/todo-api /opt/todo-api/
COPY --from=build /app/scripts /opt/todo-api/scripts
COPY --from=build /app/migrations /opt/todo-api/migrations
WORKDIR /opt/todo-api

# Run app
EXPOSE 8080
ENTRYPOINT ["./scripts/entrypoint.sh"]