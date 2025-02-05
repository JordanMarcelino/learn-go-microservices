FROM golang:1.23.4-alpine3.20 AS build-stage

WORKDIR /app

COPY . .

RUN apk add --no-cache make bash git \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && make migrateup \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -tags="sonic avx" -v -o /main ./cmd/api/main.go

FROM alpine:3.20 AS build-release-stage

WORKDIR /

RUN addgroup -S executorgroup && adduser -S executor -G executorgroup
USER executor

COPY --from=build-stage /main /main
COPY --from=build-stage /app/.env /.env

EXPOSE 8000

ENTRYPOINT [ "/main", "serve-all"]
