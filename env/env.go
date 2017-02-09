package env

// Config store the configuration of a specific environment
type Config struct {
	Address    string
	Port       int
	Restricted bool
}

// Environment is the translation map between the name of a environment and the associated configuration
var Environment map[string]*Config

func init() {
	Environment = map[string]*Config{
		"dev": &Config{
			Address:    "0.0.0.0",
			Port:       17241,
			Restricted: false,
		},
		"beta": &Config{
			Address:    "127.0.0.1",
			Port:       17242,
			Restricted: true,
		},
		"prod": &Config{
			Address:    "127.0.0.1",
			Port:       17243,
			Restricted: false,
		},
	}
}
