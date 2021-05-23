FROM golang:1.16-alpine AS build
ARG GRGATE_VERSION
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/fikaworks/grgate/pkg/config.Version=$GRGATE_VERSION'" -a -installsuffix cgo -o grgate .

FROM scratch
COPY --from=build /go/src/app/grgate /grgate
EXPOSE 8080 8086 9101
CMD ["/grgate", "serve"]
