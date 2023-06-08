FROM --platform=linux/amd64 debian:bullseye-slim

# Install Dependencies
RUN apt-get update  \
    && apt-get install -y ca-certificates wget xz-utils curl unzip python git \
    && update-ca-certificates

# Copy built binaries
RUN mkdir -p /opt/todo-api
COPY bin/todo-api /opt/todo-api/
COPY scripts /opt/todo-api/scripts
COPY migrations /opt/todo-api/migrations
WORKDIR /opt/todo-api

# Run app
EXPOSE 8080
ENTRYPOINT ["./scripts/entrypoint.sh"]