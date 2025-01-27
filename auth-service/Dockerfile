FROM golang:1.23.3-alpine3.20 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -tags="sonic avx" -v -o /main ./cmd/api/main.go

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /main /main
COPY --from=build-stage /app/.env /.env

EXPOSE 8000

ENTRYPOINT [ "/main", "serve-all"]
