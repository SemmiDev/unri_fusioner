package unri_fusioner

import "os"

type Config struct {
	SintaDomain string
	Host        string
	Port        int
}

func LoadConfig(opts ...func(c *Config)) *Config {
	config := &Config{}

	if len(opts) != 0 {
		for _, opt := range opts {
			opt(config)
		}
	}

	sintaDomain, found := os.LookupEnv("SINTA_DOMAIN")
	if found {
		config.SintaDomain = sintaDomain
	}

	host, found := os.LookupEnv("HOST")
	if found {
		config.Host = host
	}

	port, found := os.LookupEnv("PORT")
	if found {
		config.Port = CastToInt(port)
	}

	return config
}
