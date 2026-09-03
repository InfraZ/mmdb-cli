FROM --platform=${BUILDPLATFORM} golang:1.26-alpine AS build

# Provided automatically by buildx for each target in the manifest list
ARG TARGETOS
ARG TARGETARCH

# Version string baked into `mmdb-cli version`. Defaults to the in-repo placeholder
# CI passes the git tag or an "edge-<sha>" value
ARG VERSION=v0.0.0

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Mirrors .github/pre-release-hook.sh
RUN if [ "${VERSION}" != "v0.0.0" ]; then \
        sed -i "s/v0.0.0/${VERSION}/" internal/metadata/metadata.go; \
    fi

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/mmdb-cli .

FROM alpine:3

WORKDIR /data

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 1000 -h /home/infraz infraz \
    && chown infraz:infraz /data

COPY --from=build /out/mmdb-cli /usr/local/bin/mmdb-cli

LABEL org.opencontainers.image.title="mmdb-cli" \
      org.opencontainers.image.description="Command-line toolkit to create, transform, export and inspect MMDB files" \
      org.opencontainers.image.source="https://github.com/InfraZ/mmdb-cli" \
      org.opencontainers.image.documentation="https://docs.infraz.io/docs/mmdb-cli" \
      org.opencontainers.image.licenses="Apache-2.0"

USER infraz
ENTRYPOINT ["mmdb-cli"]
CMD ["--help"]
