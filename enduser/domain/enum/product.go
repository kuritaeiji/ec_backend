package enum

type ProductStatus int

const (
	OnSale = iota + 1
	SalesSuspend
	SalesEnded
)
