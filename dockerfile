# build stage
FROM golang:1.20 AS builder

ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o server

# final stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /src/server /app/server
COPY ./app.yaml /app/app.yaml
EXPOSE 8080
CMD ["/app/server"]

