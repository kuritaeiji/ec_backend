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
	"github.com/kuritaeiji/ec_backend/enduser/presentation/handler"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/middleware"
	"github.com/kuritaeiji/ec_backend/enduser/registory"
	"github.com/labstack/echo/v4"
)

// コントローラーからerror型が返却された場合のハンドラー（ログ出力し、500レスポンスを返却する）
func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(fmt.Sprintf("%+v", err))
	c.Response().Writer.WriteHeader(http.StatusInternalServerError)
}

func main() {
	e := echo.New()

	// 環境変数読み込み
	if err := config.SetupEnv(); err != nil {
		e.Logger.Fatal("環境変数読み込み失敗\n", fmt.Sprintf("%+v", err))
	}

	// コンテナ作成
	container, err := registory.NewContainer()
	if err != nil {
		e.Logger.Fatal("コンテナ作成失敗\n", fmt.Sprintf("%+v", err))
	}

	e, loginG, err := middleware.SetupMiddleware(e, container)
	if err != nil {
		e.Logger.Fatal("ミドルウェア設定失敗\n", fmt.Sprintf("%+v", err))
	}

	// エラーハンドラー設定
	e.HTTPErrorHandler = customHTTPErrorHandler

	// ハンドラー設定
	err = handler.SetupHandlers(e, loginG, container)
	if err != nil {
		e.Logger.Fatal("ハンドラー設定失敗\n", fmt.Sprintf("+%v", err))
	}

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
