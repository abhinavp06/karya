FROM golang:1.26-alpine AS build
RUN apk add --no-cache build-base
WORKDIR /src
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /karya .

FROM alpine:3.22
COPY --from=build /karya /karya
EXPOSE 3000
ENTRYPOINT ["/karya"]