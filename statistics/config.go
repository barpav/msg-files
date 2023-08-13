package statistics

import "os"

const (
	defaultHost     = "localhost"
	defaultPort     = "5672"
	defaultUser     = "guest"
	defaultPassword = "guest"
)

const (
	envVarHost     = "MSG_FILES_STAT_HOST"
	envVarPort     = "MSG_FILES_STAT_PORT"
	envVarUser     = "MSG_FILES_STAT_USER"
	envVarPassword = "MSG_FILES_STAT_PASSWORD"
)

type config struct {
	host     string
	port     string
	user     string
	password string
}

func (c *config) read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readSetting(envVarPort, defaultUser, &c.user)
	readSetting(envVarPort, defaultPassword, &c.password)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}
