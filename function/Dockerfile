# syntax=docker/dockerfile:1

ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.22

################################################################################
# Create a stage for building the application.
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
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
FROM alpine:${ALPINE_VERSION} AS final

ARG APP_PORT

# Install any runtime dependencies that are needed to run your application.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        curl \
        && \
        update-ca-certificates

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/main /bin/

# Expose the port that the application listens on.
ENV FUNCTIONS_CUSTOMHANDLER_PORT=${APP_PORT}
EXPOSE ${APP_PORT}

# What the container should run when it is started.
ENTRYPOINT [ "/bin/main" ]


################################################################################
# Create a new stage for building the image specific for azure app service.
FROM alpine:${ALPINE_VERSION} AS appservice

ARG APP_PORT=8080

# Copy app service ssh config
COPY ./appservice/sshd_config /etc/ssh/sshd_config
COPY ./appservice/entrypoint.sh /

# Start and enable SSH
RUN apk add openssh \
    && echo "root:Docker!" | chpasswd \
    && chmod +x /entrypoint.sh \
    && cd /etc/ssh/ \
    && ssh-keygen -A

# Copy the executable from the "build" stage.
COPY --from=build /bin/main /bin/

# Expose the port that the application listens on.
ENV APP_PORT=${APP_PORT}
EXPOSE 2222 ${APP_PORT}
