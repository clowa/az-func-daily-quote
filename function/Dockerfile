# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25

################################################################################
# Create a stage for building the application.
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS builder
WORKDIR /app

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Copy the application source code into the container.
COPY . .

# This is the architecture you’re building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
WORKDIR /app/src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/main 

################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application. This often uses a different base
# image from the build stage where the necessary files are copied from the build
# stage.
FROM scratch

ARG APP_PORT

# Copy TLS certificates to allow TLS traffic
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the executable from the "build" stage.
COPY --from=builder /bin/main /bin/

# Expose the port that the application listens on.
ENV FUNCTIONS_CUSTOMHANDLER_PORT=${APP_PORT}
EXPOSE ${APP_PORT}

# What the container should run when it is started.
ENTRYPOINT [ "/bin/main" ]
