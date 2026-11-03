FROM golang:{{ cookiecutter.go_version }}-alpine AS build

RUN apk add --no-cache git

# This argument should be defined as your personal access tokens, from: https://github.com/settings/tokens
ARG GITHUB_ACCESS_TOKEN

RUN git config \
    --global \
    url."https://${GITHUB_ACCESS_TOKEN}:@github.com".insteadOf \
    "https://github.com"

WORKDIR /src
RUN mkdir -p /src/api/go

# Caching layers for dependencies
COPY ./go.mod ./go.sum* ./
COPY ./api/go/go.mod ./api/go/go.sum ./api/go/
COPY ./go.work ./go.work.sum* ./

RUN go env -w GOMODCACHE=/go/pkg/mod GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN mkdir -p /app/bin

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GIT_TERMINAL_PROMPT=1 \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o /app/bin ./...

# https://github.com/GoogleContainerTools/distroless/blob/master/examples/go/Dockerfile
# Use distroless would loss us the ability to k exec into the pod
FROM gcr.io/distroless/cc-debian12
COPY --from=build /app/bin/* /

ARG GIT_SHA=unknown
ARG BUILD_TIME=unknown
ENV GIT_SHA=$GIT_SHA
ENV BUILD_TIME=$BUILD_TIME

EXPOSE 8080 9090 4014

ENTRYPOINT ["/server"]
