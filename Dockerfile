FROM golang:1.20-alpine AS build
ARG GRGATE_COMMITSHA
ARG GRGATE_VERSION
RUN apk --no-cache add ca-certificates
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X 'github.com/fikaworks/grgate/pkg/config.Version=$GRGATE_VERSION' -X 'github.com/fikaworks/grgate/pkg/config.CommitSha=$GRGATE_COMMITSHA'" \
  -a -installsuffix cgo -o grgate .

FROM scratch
LABEL org.opencontainers.image.source https://github.com/FikaWorks/grgate
COPY --from=build /go/src/app/grgate /bin/grgate
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080 8086 9101
CMD ["grgate", "serve"]
