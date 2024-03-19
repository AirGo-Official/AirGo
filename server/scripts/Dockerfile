FROM --platform=$BUILDPLATFORM golang:alpine AS build-server
WORKDIR /server
COPY server /server
ARG TARGETOS
ARG TARGETARCH
RUN set -ex \
    && GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o AirGo -trimpath -ldflags "-s -w" main.go

FROM --platform=${TARGETPLATFORM} alpine:latest
WORKDIR /air
COPY --from=build-server /server/AirGo ./
COPY --from=build-server /server/config.yaml ./

VOLUME ["/air"]

ENV LANG C.UTF-8
ENTRYPOINT [ "/air/AirGo" ]
CMD ["start"]