# syntax=docker/dockerfile:1

FROM golang:1.19.2-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY . ./

ENV DATA_DIR=/app/resources
ENV LOG_DIR=/

WORKDIR /app/main
RUN go build -o /docker-servicediscovery

EXPOSE 8000

CMD [ "/docker-servicediscovery" ]