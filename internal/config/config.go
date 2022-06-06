package config

type Config struct {
	Address        string
	DSN            string
	AccrualAddress string
	Secret         string
}

type Application struct {
	Config   Config
	Sessions map[string]string
}
