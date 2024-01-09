package config

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func SetupLogger(e *echo.Echo) {
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} [${level}] ${long_file}:${line}\n")
		levStr := os.Getenv("LOG_LEVEL")
		levInt, err := strconv.Atoi(levStr)
		if err != nil {
			return
		}

		l.SetLevel(log.Lvl(levInt))
	}
}
