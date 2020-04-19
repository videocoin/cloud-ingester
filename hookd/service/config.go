package service

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	Addr           string `envconfig:"ADDR" default:"0.0.0.0:8887"`
	StreamsRPCAddr string `envconfig:"STREAMS_RPC_ADDR" default:"127.0.0.1:5102"`
	EmitterRPCAddr string `envconfig:"EMITTER_RPC_ADDR" default:"127.0.0.1:5003"`
	HLSDir         string `default:"/tmp/hls"`
}
