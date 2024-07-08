# Stage 1: Base
FROM golang:1.22-alpine AS base
RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Stage 2: Build
FROM base AS build

ENV PORT = 443
RUN go build -o /bin/main ./cmd/web/main.go

# Stage 3: Final image
FROM alpine:latest
# Copy resources
COPY  resources /resources
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/main /usr/local/bin/main

COPY certificates /certificates
COPY . .

EXPOSE 443

ENTRYPOINT ["main"]