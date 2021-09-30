package supermarket

type SpecialOfferType int

const (
	TenPercentDiscount SpecialOfferType = iota
	ThreeForTwo
	TwoForAmount
	FiveForAmount
)

type SpecialOffer struct {
	offerType SpecialOfferType
	product           Product
	costAfterDiscount float64
}

type Discount struct {
	product Product
	description string
	discountAmount float64
}

