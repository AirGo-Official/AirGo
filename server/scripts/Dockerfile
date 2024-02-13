#FROM --platform=$BUILDPLATFORM ubuntu:latest AS build-web
#WORKDIR /web
#COPY /web /web
#RUN set -ex \
#    && apt update -y && apt install npm curl wget tar unzip -y \
#    && npm install vite -g && npm install -g n && n lts \
#    && npm -v && node -v \
#    && npm install && npm run build

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

ENV TZ=Asia/Shanghai
ENV LANG C.UTF-8
ENTRYPOINT [ "/air/AirGo" ]
CMD ["-start"]