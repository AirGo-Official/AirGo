FROM --platform=$TARGETPLATFORM alpine:latest
ARG TARGETOS
ARG TARGETARCH
#ARG BUILDPLATFORM
#ARG TARGETPLATFORM

#linux/arm,linux/arm64,linux/amd64

WORKDIR /air

VOLUME ["/air/config.yaml"]

COPY server/AirGo-$TARGETOS-$TARGETARCH /air/AirGo

ENV LANG C.UTF-8

EXPOSE 80 443

CMD ["/air/AirGo", "-start"]

# docker buildx build --platform linux/arm,linux/arm64,linux/amd64 -f server/scripts/Dockerfile -t airgo .