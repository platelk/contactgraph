## Build containers
FROM  golang:1.17.5-alpine3.15 AS build

RUN apk update && apk add make git gcc musl-dev

WORKDIR /app
COPY . .

RUN go build -o bin

## Running containers
FROM alpine:3.15
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/bin /app/bin

RUN chmod +x /app/bin

ENTRYPOINT exec /app/bin $0 $@
