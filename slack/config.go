package slack

type Config struct {
	// Whether Slack integration is enabled.
	Enabled bool `toml:"enabled"`
	// The Slack webhook URL, can be obtained by adding Incoming Webhook integration.
	TokenKey string `toml:"tokenkey"`
	// The default channel, can be overridden per alert.
}

func NewConfig() Config {
	return Config{}
}
