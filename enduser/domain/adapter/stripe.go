//go:generate mockery --name StripeAdapter
package adapter

type StripeAdapter interface {
	CreateCustomer() (string, error)
}
