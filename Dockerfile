# The context is assembled by goreleaser or by `make image` and already holds the compiled
# binary and the license, so this copies rather than builds.
FROM scratch

# dockers_v2 lays the build context out as linux/<arch>/<binary>, so the binary is not at
# the context root. `make image` mirrors that layout for this reason.
ARG TARGETPLATFORM

# The controller presents a certificate, so the trust store has to be present even
# though nothing else is. It comes from a pinned image rather than the build context,
# which is what lets `make image` and goreleaser share one Dockerfile.
COPY --from=alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY LICENSE NOTICE /
COPY $TARGETPLATFORM/wnc /wnc

USER 65534:65534

ENTRYPOINT ["/wnc"]
CMD ["--help"]

LABEL org.opencontainers.image.title="wnc"
LABEL org.opencontainers.image.description="CLI for Cisco Catalyst 9800 Wireless Network Controllers"
LABEL org.opencontainers.image.vendor="umatare5"
LABEL org.opencontainers.image.source="https://github.com/umatare5/wnc"
LABEL org.opencontainers.image.licenses="MIT"
