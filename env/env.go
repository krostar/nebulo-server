package env

// Config represent an environment with differents kind of informations
type Config struct {
	Address string
	Port    int
}

// Environment is the type of all the available environment
type Environment string

const (
	// DEV is the environment used in development,
	// with more details in case of HTTP errors
	DEV = Environment("dev")

	// PREPROD is the environment used in pre-production,
	// the only difference with production environment is the services used
	PREPROD = Environment("preprod")

	// PROD is the environment used in production,
	// with inexistant debug in case of error client-side
	PROD = Environment("prod")
)

var (
	// EnvironmentConfig store the configuration of each environments
	EnvironmentConfig map[Environment]*Config
)

func init() {
	EnvironmentConfig = map[Environment]*Config{
		DEV: &Config{
			Address: "0.0.0.0",
			Port:    17241,
		},
		PREPROD: &Config{
			Address: "127.0.0.1",
			Port:    17242,
		},
		PROD: &Config{
			Address: "127.0.0.1",
			Port:    17243,
		},
	}
}
