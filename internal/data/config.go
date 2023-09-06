package data

import "os"

const (
	defaultHost = "localhost"
	defaultPort = "27017"
)

const (
	envVarHost = "MSG_STORAGE_HOST"
	envVarPort = "MSG_STORAGE_PORT"
)

type config struct {
	host string
	port string
}

func (c *config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}
