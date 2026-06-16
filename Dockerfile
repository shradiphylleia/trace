FROM golang:1.22-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o /out/traceshare ./cmd/server

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /out/traceshare /app/traceshare
EXPOSE 8080
ENTRYPOINT ["/app/traceshare"]
