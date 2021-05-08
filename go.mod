module github.com/fikaworks/ggate

go 1.16

require (
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/golang/mock v1.3.1
	github.com/google/go-github/v34 v34.0.0
	github.com/gorilla/mux v1.8.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/justinas/alice v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/rs/zerolog v1.21.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/xanzy/go-gitlab v0.49.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
)

replace github.com/xanzy/go-gitlab => github.com/etiennetremel/go-gitlab v0.49.1-0.20210506152720-85530398e40b
