# syntax=docker/dockerfile:1

# Keep the container toolchain aligned with Mise. Renovate updates both.

# ---- Go build -------------------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.27.2-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=unknown
ARG BRANCH=unknown
ARG BUILD_DATE=unknown

RUN apk add --no-cache upx
WORKDIR /workspace

# Cache module downloads before copying source.
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/google/go-licenses/v2@v2.0.1

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go-licenses save ./cmd/selectronic_exporter --save_path third_party_licenses --ignore github.com/woodleighschool/selectronic-exporter --force

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath \
    -ldflags "-s -w -X github.com/prometheus/common/version.Version=${VERSION#v} -X github.com/prometheus/common/version.Revision=${COMMIT} -X github.com/prometheus/common/version.Branch=${BRANCH} -X github.com/prometheus/common/version.BuildUser=docker -X github.com/prometheus/common/version.BuildDate=${BUILD_DATE}" \
    -o selectronic_exporter ./cmd/selectronic_exporter
RUN upx --best --lzma selectronic_exporter

# ---- Runtime --------------------------------------------------------------
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY LICENSE /LICENSE
COPY --from=builder /workspace/third_party_licenses /third_party_licenses
COPY --from=builder /usr/local/go/LICENSE /third_party_licenses/go/LICENSE
COPY --from=builder /workspace/selectronic_exporter /selectronic_exporter
EXPOSE 9788
USER 65532:65532
ENTRYPOINT ["/selectronic_exporter"]
