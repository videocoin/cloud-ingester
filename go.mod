module github.com/videocoin/cloud-ingester

go 1.12

require (
	github.com/evalphobia/logrus_sentry v0.8.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/videocoin/cloud-api v0.2.7
	github.com/videocoin/hookd v0.0.66
	google.golang.org/grpc v1.23.0
)

replace github.com/videocoin/cloud-api => ../cloud-api
