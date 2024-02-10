//go:generate mockery --name EmailAdapter
package adapter

type EmailAdapter interface {
	SendEmail(from string, to string, subject string, text string) error
}
