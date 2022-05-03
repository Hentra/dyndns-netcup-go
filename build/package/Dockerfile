# build stage
FROM golang:1.18.1-alpine3.15 as build-stage
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /dyndns-docker ./cmd/dyndns-netcup-docker/main.go 

# production stage
LABEL org.opencontainers.image.source https://github.com/Hentra/dyndns-netcup-go
FROM alpine:latest
RUN mkdir /cache/
COPY --from=build-stage /dyndns-docker /dyndns-docker

ENTRYPOINT ["/dyndns-docker"]

