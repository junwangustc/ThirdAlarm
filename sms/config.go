package sms

type Config struct {
	Enabled   bool   `toml:"enabled"`
	SenderKey string `toml:"senderkey"`
	Token     string `toml:"token"`
	SendUrl   string `toml:"sendurl"`
	Expandkey string `toml:"expandkey"`
}

func NewConfig() Config {
	return Config{}
}
