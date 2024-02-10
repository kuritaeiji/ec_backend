package bridge

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/cockroachdb/errors"
)

type emailAdapter struct{}

const (
	From = "info@ec-site.shop"
)

func NewEmailAdapter() emailAdapter {
	return emailAdapter{}
}

// メールを送信する
func (ea emailAdapter) SendEmail(from string, to string, subject string, text string) error {
	// AWSセッションの作成
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		return err
	}

	// SESクライアントの作成
	svc := ses.New(sess)

	// メール送信パラメータの作成
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(text),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	// メールを送信
	_, err = svc.SendEmail(params)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
