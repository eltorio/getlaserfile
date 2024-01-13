FROM golang:1.21 as builder
COPY src /src
WORKDIR /src
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o getlaserfile .

FROM ubuntu:latest
COPY --from=builder /src/getlaserfile /usr/local/bin/getlaserfile
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]