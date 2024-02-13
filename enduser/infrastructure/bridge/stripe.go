package bridge

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
)

type stripeAdapter struct{
	secretKey string
}

func NewStripeAdapter() stripeAdapter {
	return stripeAdapter{
		secretKey: os.Getenv("STRIPE_SECRET_KEY"),
	}
}

// StripeにCustomerを作成する。返り値として作成したCustomerのIDを受け取る。
func (sa stripeAdapter) CreateCustomer() (string, error) {
	stripe.Key = sa.secretKey

	params := &stripe.CustomerParams{}
	result, err := customer.New(params)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return result.ID, nil
}
