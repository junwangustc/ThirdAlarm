package lark

type Config struct {
	// Whether lark integration is enabled.
	Enabled bool `toml:"enabled"`
	// The lark webhook URL, can be obtained by adding Incoming Webhook integration.
	UserTokenKey string `toml:"usertokenkey"`
	// The default channel, can be overridden per alert.
	BotTokenKey string `toml:"bottokenkey"`
	ChannelID   string `toml:"channelid"`
	Url         string `toml:"url"`
}

func NewConfig() Config {
	return Config{}
}
