package test

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/labstack/gommon/log"
)

// フォームオブジェクトをio.Reader型に変換する
func FormToReader(form any) io.Reader {
	// JSONエンコード
	jsonData, err := json.Marshal(form)
	if err != nil {
		log.Errorf("%+vをreaderに変換できませんでした", form)
		return nil
	}

	// JSONデータをio.Readerに変換
	reader := bytes.NewReader(jsonData)
	return reader
}

// io.Readerを引数responseに変換する
func ReaderToResponse(target io.Reader, response any) {
	err := json.NewDecoder(target).Decode(response)
	if err != nil {
		log.Errorf("%+vをフォームに変換できませんでした", response)
	}
}

// 引数に渡された変数のポインターを返却する
func ToPointer[T any](target T) *T {
	return &target
}