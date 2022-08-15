package config

type config interface {
	ListenPort() string
	ListenAddress() string
}

var cfg config

func SetConfig(c config) {
	cfg = c
}

func ListenPort() string {
	return cfg.ListenPort()
}

func ListenAddress() string {
	return cfg.ListenAddress()
}
