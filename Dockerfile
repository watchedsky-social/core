FROM cgr.dev/chainguard/go:latest AS build

ARG VERSION=dev
ARG WATCHEDSKY_DB_PASSWORD

ENV GOOS=linux
ENV GOARCH=amd64
ENV WATCHEDSKY_ENV_FILE=/run/secrets/env

WORKDIR /src
COPY . /src/

RUN ["mkdir", "/assets"]
RUN ["go", "mod", "download"]
RUN --mount=type=secret,id=env ["go", "generate", "./..."]
RUN ["go", "build", "-trimpath", "-ldflags", "-X main.Version=${VERSION}", "-o", "/assets/core", "main.go"]

FROM cgr.dev/chainguard/glibc-dynamic:latest AS release

WORKDIR /app

COPY --from=build /assets/core .

RUN ["/app/core", "install"]
