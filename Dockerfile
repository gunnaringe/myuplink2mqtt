FROM golang:alpine AS build
RUN apk add git
RUN apk add --no-cache $(apk search --no-cache | grep -q ^upx && echo -n upx)

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v ./...
RUN env CGO_ENABLED=0 go build -ldflags '-w -s' -o /go/bin/app cmd/main.go

RUN if command -v upx &> /dev/null; then upx --ultra-brute /go/bin/app; fi

FROM gcr.io/distroless/static
LABEL maintainer="Gunnar Inge G. Sortland <gunnar.inge@sort.land>"
COPY --from=build /go/bin/app /
ENTRYPOINT [ "/app"]
