package env

// Config represent an environment with differents kind of informations
type Config struct {
	Type    string `json:"type" validate:"regexp=^(dev|preprod|prod)$"`
	Address string `json:"address,omitempty" validate:"-"`
	Port    int    `json:"port,omitempty" validate:"-"`
}

const (
	// DEV is the environment used in development,
	// with more details in case of HTTP errors
	DEV = "dev"

	// PREPROD is the environment used in pre-production,
	// the only difference with production environment is the services used
	PREPROD = "preprod"

	// PROD is the environment used in production,
	// with inexistant debug in case of error client-side
	PROD = "prod"
)

var (
	// EnvironmentConfig store the configuration of each environments
	EnvironmentConfig map[string]*Config
)

func init() {
	EnvironmentConfig = map[string]*Config{
		DEV: &Config{
			Type:    DEV,
			Address: "0.0.0.0",
			Port:    17241,
		},
		PREPROD: &Config{
			Type:    PREPROD,
			Address: "127.0.0.1",
			Port:    17242,
		},
		PROD: &Config{
			Type:    PROD,
			Address: "127.0.0.1",
			Port:    17243,
		},
	}
}
