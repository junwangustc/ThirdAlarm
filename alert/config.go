package alert

type URL struct {
	Name      string   `toml:"name"`
	Url       string   `toml:"url"`
	CareSlack []string `toml:"careslack"`
	CareSMS   []string `toml:"caresms"`
}
type Config struct {
	LiveUrl []URL `toml:"liveurl"`
}

func NewConfig() Config {
	return Config{}
}
