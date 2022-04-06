ARG gover=1.16.4

FROM golang:$gover as build

WORKDIR /opt/app

COPY . .

RUN CGO_ENABLED=0 go build -o imgen -a -installsuffix cgo

FROM alpine:3.15.4

WORKDIR /opt/app

COPY --from=build /opt/app/imgen .

CMD ["./imgen"]
