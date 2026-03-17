# set app name
ARG app=pinglog

# create build stage
ARG TAG=alpine
FROM --platform=$BUILDPLATFORM golang:$TAG AS build
ARG app

# copy source files into the container
COPY . /src/$app/

# build, strip, and compress the binary
WORKDIR /src/$app
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    GOEXPERIMENT=goroutineleakprofile \
    go build -trimpath -ldflags "-s -w" -tags "netgo timetzdata" -o $app \
    && chmod 500 $app

# set up final stage
FROM scratch
ARG app

# copy in user info
COPY --chown=root:root --chmod=0400 docker/passwd /etc/passwd

# run as nonroot
USER nonroot

# copy in binary
COPY --from=build --chown=root:root --chmod=0005 /src/$app/$app /$app

# run application
ENTRYPOINT ["/pinglog"]
