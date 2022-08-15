package config

type config interface {
	ListenPort() string
	ListenAddress() string
	PiggyBankDBHost() string
	PiggyBankDBPort() string
	PiggyBankDBName() string
	PiggyBankDBUsername() string
	PiggyBankDBPassword() string
	PiggyBankDBSSLMode() string
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

func PiggyBankDBHost() string {
	return cfg.PiggyBankDBHost()
}

func PiggyBankDBPort() string {
	return cfg.PiggyBankDBPort()
}

func PiggyBankDBName() string {
	return cfg.PiggyBankDBName()
}

func PiggyBankDBUsername() string {
	return cfg.PiggyBankDBUsername()
}

func PiggyBankDBPassword() string {
	return cfg.PiggyBankDBPassword()
}

func PiggyBankDBSSLMode() string {
	return cfg.PiggyBankDBSSLMode()
}
