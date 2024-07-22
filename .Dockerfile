# Stage 1: Base
FROM golang:1.22.3-alpine AS base
RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
ARG PORT=443

COPY aggregator/ ./aggregator/
COPY client/ ./client/
COPY cmd/ ./cmd/
COPY collector/ ./collector/
COPY constant/ ./constant/
COPY entity/ ./entity/
COPY filter/ ./filter/
COPY parser/ ./parser/
COPY resources/ ./resources/
COPY sorter/ ./sorter/
COPY validator/ ./validator/
COPY storage/ ./storage/

RUN go build -o /bin/main ./cmd/web/main.go

# Stage 2: Build image
FROM alpine:3.20.1
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /bin/main /usr/local/bin/main

# Copy storage to the root directory
COPY --from=base /src/storage /storage
# Copy resources to the root directory
COPY --from=base /src/resources /resources

COPY certificates /certificates

EXPOSE ${PORT}

ENTRYPOINT ["main"]
