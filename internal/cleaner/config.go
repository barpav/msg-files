package cleaner

import (
	"os"
	"strconv"
)

const (
	defaultCheckPeriod = 300 // seconds
	defaultDeleteAfter = 300 // seconds
)

const (
	envVarCheckPeriod = "MSG_UNUSED_FILES_CHECK_PERIOD"
	envVarDeleteAfter = "MSG_UNUSED_FILES_DELETE_AFTER"
)

type config struct {
	checkPeriod int
	deleteAfter int
}

func (c *config) read() {
	readNumericSetting(envVarCheckPeriod, defaultCheckPeriod, &c.checkPeriod)

	if c.checkPeriod < 10 {
		c.checkPeriod = defaultCheckPeriod
	}

	readNumericSetting(envVarDeleteAfter, defaultDeleteAfter, &c.deleteAfter)

	if c.checkPeriod < 10 {
		c.checkPeriod = defaultCheckPeriod
	}
}

func readNumericSetting(setting string, defaultValue int, result *int) {
	val := os.Getenv(setting)

	if val != "" {
		valNum, err := strconv.Atoi(val)

		if err == nil {
			*result = valNum
			return
		}
	}

	*result = defaultValue
}
