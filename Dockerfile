FROM docker.io/golang:1.22.2@sha256:d5302d40dc5fbbf38ec472d1848a9d2391a13f93293a6a5b0b87c99dc0eaa6ae as builder
ARG VERSION="devel"

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.version=${VERSION}'" -o kie ./cmd/kubernetes-image-exporter

FROM docker.io/alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

# renovate: datasource=repology depName=alpine_3_19/ca-certificates versioning=loose
ENV CA_CERTIFICATES_VERSION="20240226-r0"
# renovate: datasource=repology depName=alpine_3_19/dumb-init versioning=loose
ENV DUMB_INIT_VERSION="1.2.5-r3"

COPY --from=builder /build/kie /kie

RUN apk add --no-cache \
    ca-certificates="${CA_CERTIFICATES_VERSION}" \
    dumb-init="${DUMB_INIT_VERSION}"

EXPOSE 9145
USER 65534

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/kie", "serve"]
