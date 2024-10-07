FROM golang:1.22 AS build

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o imgen -a -installsuffix cgo

FROM alpine:3.20

RUN addgroup -S imgen && adduser -S imgen -G imgen

WORKDIR /opt/app

COPY --from=build /opt/app/imgen .

RUN chown imgen:imgen imgen
USER imgen

CMD ["./imgen"]
