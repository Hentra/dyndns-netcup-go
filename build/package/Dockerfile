# build stage
FROM golang:1.24.2-alpine3.21 AS build-stage
WORKDIR /app
COPY . .
RUN touch /ipcache
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata -o /dyndns-docker ./cmd/dyndns-netcup-docker/main.go

# production stage
LABEL org.opencontainers.image.source=https://github.com/Hentra/dyndns-netcup-go
FROM scratch
USER 62534:62534
COPY --from=build-stage --chown=62534:62534 /ipcache /ipcache
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /dyndns-docker /dyndns-docker

ENTRYPOINT ["/dyndns-docker"]

