module github.com/videocoin/cloud-ingester

go 1.12

require (
	github.com/grafov/m3u8 v0.11.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/opentracing-contrib/echo v0.0.0-20190807091611-5fe2e1308f06
	github.com/opentracing/opentracing-go v1.1.0
	github.com/plutov/echo-logrus v1.0.0
	github.com/prometheus/client_golang v1.4.1
	github.com/sirupsen/logrus v1.5.0
	github.com/videocoin/cloud-api v0.2.7
	github.com/videocoin/cloud-pkg v0.0.5
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg
