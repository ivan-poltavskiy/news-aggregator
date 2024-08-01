FROM golang:1.22.3-alpine AS base

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY news-updater/go.mod news-updater/go.sum ./

COPY web/ ./web/
COPY entity/ ./entity/
COPY constant/ ./constant/
COPY storage/ ./storage/
COPY parser/ ./parser/
COPY news-updater/ ./news-updater/

WORKDIR /src/news-updater
RUN go mod download

RUN go build -o /bin/news-updater ./main.go

FROM alpine:3.20.1
COPY --from=base /bin/news-updater .

ENTRYPOINT ["/news-updater"]
