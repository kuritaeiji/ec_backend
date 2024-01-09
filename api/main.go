package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kuritaeiji/ec_backend/config"
	"github.com/kuritaeiji/ec_backend/presentation/handler"
	"github.com/labstack/echo/v4"
)

// コントローラーからerror型が返却された場合のハンドラー（ログ出力し、500レスポンスを返却する）
func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(fmt.Sprintf("%+v", err))
	c.Response().Writer.WriteHeader(http.StatusInternalServerError)
}

// ハンドラーのセットアップ
func setupHandlers(e *echo.Echo) {
	handler.SetupHealthcheckHandler(e)
}

func main() {
	e := echo.New()

	// 環境変数読み込み
	if err := config.SetupEnv(); err != nil {
		e.Logger.Error("環境変数読み込み失敗\n", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	// DB接続
	_, close, err := config.SetupDB()
	if err != nil {
		e.Logger.Error("DB接続失敗\n", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
	defer close()

	// ログ設定
	config.SetupLogger(e)

	// エラーハンドラー設定
	e.HTTPErrorHandler = customHTTPErrorHandler

	// ハンドラー設定
	setupHandlers(e)

	// サーバー起動（FargateのSIGINTシグナルを受け取ると停止するようにする）
	errorCh := make(chan error, 1)
	go func() {
		err := e.Start(":8080")
		if err != nil {
			errorCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errorCh:
		e.Logger.Error("サーバー起動失敗\n", fmt.Sprintf("%+v", err))
		os.Exit(1)
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Error("サーバーシャットダウン失敗\n", fmt.Sprintf("%+v", err))
			os.Exit(1)
		} else {
			e.Logger.Info("サーバーシャットダウン")
		}
	}
}