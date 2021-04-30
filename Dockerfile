FROM golang:1.16-alpine AS build
ARG GGATE_VERSION
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/fikaworks/ggate/pkg/config.Version=$GGATE_VERSION'" -a -installsuffix cgo -o ggate .

FROM scratch
COPY --from=build /go/src/app/ggate /ggate
EXPOSE 8080 8086 9101
CMD ["/ggate", "serve"]
