FROM golang:1.16 AS build-env
ADD . /go/src/github.com/grpc-streamer
WORKDIR /go/src/github.com/grpc-streamer
 
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' cmd/main.go

FROM alpine as cert
RUN apk update && apk add ca-certificates


FROM alpine:3.8
COPY --from=build-env /go/src/github.com/grpc-streamer /usr/local/bin/server/
COPY --from=cert /etc/ssl/certs /etc/ssl/certs
ENV ENV="local"

ENV GRPC_HEALTH_PROBE_VERSION v0.3.2
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \chmod +x /bin/grpc_health_probe

EXPOSE 50052
WORKDIR /usr/local/bin/server
CMD ./main
