# syntax=docker/dockerfile:1

# Defaults keep local builds self-contained. Renovate updates this alongside
# the matching Mise and module pins.
ARG GO_VERSION=1.26.5

# ---- Go build -------------------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache upx
WORKDIR /workspace

# Cache module downloads before copying source.
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w" -o selectronic_exporter ./cmd/selectronic_exporter
RUN upx --best --lzma selectronic_exporter

# ---- Runtime --------------------------------------------------------------
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/selectronic_exporter /selectronic_exporter
EXPOSE 9788
USER 65532:65532
ENTRYPOINT ["/selectronic_exporter"]
