FROM golang:latest AS build
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

RUN apt-get install git

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v ./...
RUN env CGO_ENABLED=0 go build -ldflags '-w -s' -o /go/bin/app cmd/main.go

FROM scratch
LABEL maintainer="Gunnar Inge G. Sortland <gunnar.inge@sort.land>"
LABEL org.opencontainers.image.description="Brigde for myUplink and Home Assistant"
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc_passwd /etc/passwd
USER nobody

COPY --from=build /go/bin/app /
ENTRYPOINT [ "/app"]
