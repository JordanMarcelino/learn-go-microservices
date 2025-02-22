FROM scratch AS build-stage

WORKDIR /app

COPY . .

FROM alpine:3.20 AS build-release-stage

WORKDIR /

RUN addgroup -S executorgroup && adduser -S executor -G executorgroup
USER executor

COPY --from=build-stage /app/bin/main /main
COPY --from=build-stage /app/.env /.env

EXPOSE 8000

ENTRYPOINT [ "/main", "serve-all" ]
