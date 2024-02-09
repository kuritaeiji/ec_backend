package util

// バリデーションエラー時に表示されるフィールド名
var FieldNames = map[string]string{
	"ValidationAccountForCreation.Email":                "メールアドレス",
	"ValidationAccountForCreation.Password":             "パスワード",
	"ValidationAccountForCreation.PasswordConfirmation": "パスワード（確認用）",
	"ValidationAccountForReviewNickname.ReviewNickname": "レビュー投稿者名",
}
