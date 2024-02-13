package util

import (
	"log"
	"time"
)

type (
	TimeUtils interface {
		NowJP() time.Time
		DateJP(year int, month time.Month, day int) time.Time
		DateToUTC(t time.Time) time.Time
		DateToJP(t time.Time) time.Time
		TimeToUTC(t time.Time) time.Time
		TimeToJP(t time.Time) time.Time
	}
	timeUtils struct {
		loc *time.Location
	}
)

func NewTimeUtils() timeUtils {
	// 日本時間のタイムゾーンを取得
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	return timeUtils{loc: loc}
}

// 日本のタイムゾーンの現在日時を返却する
func (tu timeUtils) NowJP() time.Time {
	return time.Now().In(tu.loc)
}

// 日本のタイムゾーンの日時を返却する
func (tu timeUtils) DateJP(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, tu.loc)
}

// ゾーンをUTCから日本時間に変換する
// 日付の場合に使用する
func (tu timeUtils) DateToUTC(t time.Time) time.Time {
	return t.In(tu.loc)
}

// ゾーンを日本からUTCに変換する
// 日付の場合に使用する
func (tu timeUtils) DateToJP(t time.Time) time.Time {
	return t.In(time.UTC)
}

// 日本時間からUTC時間に変換する
func (tu timeUtils) TimeToUTC(t time.Time) time.Time {
	return t.Add(-9 * time.Hour).In(time.UTC)
}

// UTC時間から日本時間に変換する
func (tu timeUtils) TimeToJP(t time.Time) time.Time {
	return t.Add(9 * time.Hour).In(tu.loc)
}
