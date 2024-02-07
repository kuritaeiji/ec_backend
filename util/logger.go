package util

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func NewLogger() echo.Logger {
	e := echo.New()
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} [${level}] ${long_file}:${line}\n")
		levStr := os.Getenv("LOG_LEVEL")
		levInt, err := strconv.Atoi(levStr)
		if err != nil {
			e.Logger.Error("LOG_LEVEL環境変数が正しく設定されていません\nログレベルをエラーにしました")
			levInt = int(log.ERROR)
		}

		l.SetLevel(log.Lvl(levInt))
	}

	return e.Logger
}
