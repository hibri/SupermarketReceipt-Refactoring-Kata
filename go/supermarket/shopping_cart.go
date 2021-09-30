package supermarket

import (
	"fmt"
	"math"
)

type ProductQuantity struct {
	product Product
	quantity float64
}

type ShoppingCart struct {
	items []ProductQuantity
	productQuantities map[Product]float64
}

func NewShoppingCart() *ShoppingCart {
	var s ShoppingCart
	s.items = []ProductQuantity{}
	s.productQuantities = make(map[Product]float64)
	return &s
}

func (c *ShoppingCart) addItem(product Product) {
	c.addItemQuantity(product, 1)
}

func (c *ShoppingCart) addItemQuantity(product Product, amount float64) {
	c.items = append(c.items, ProductQuantity{product: product, quantity: amount})
	currentAmount, ok := c.productQuantities[product]
	if ok {
		c.productQuantities[product] = currentAmount + amount
	} else {
		c.productQuantities[product] = amount
	}
}

func (c *ShoppingCart) handleOffers(receipt *Receipt, offers map[Product]SpecialOffer, catalog Catalog) {
	for product, _ := range c.productQuantities {
		var quantity = c.productQuantities[product]
		offer, offerExists := offers[product]
		if offerExists  {
			applyOffer(receipt, catalog, product, quantity, offer)
		}

	}
}

func applyOffer(receipt *Receipt, catalog Catalog, product Product, quantity float64, offer SpecialOffer) {
	var unitPrice = catalog.unitPrice(product)
	var quantityAsInt = int(math.Round(quantity))
	var discount *Discount = nil
	var x = 1
	switch offer.offerType {
	case ThreeForTwo:
		x = 3
		break
	case TwoForAmount:
		x = 2
		break
	case FiveForAmount:
		x = 5
		break
	default:
		x = 1
	}
	if offer.offerType == ThreeForTwo {
		x = 3

	} else if offer.offerType == TwoForAmount {
		x = 2
		if quantityAsInt >= 2 {
			discount = calculateDiscountForTwo(offer, quantityAsInt, x, unitPrice, quantity, discount, product)
		}

	}
	if offer.offerType == FiveForAmount {
		x = 5
	}
	var numberOfXs int = quantityAsInt / x
	if offer.offerType == ThreeForTwo && quantityAsInt > 2 {
		discount = calculateDiscountThreeForTwo(quantity, unitPrice, numberOfXs, quantityAsInt, discount, product)
	}
	if offer.offerType == TenPercentDiscount {
		discount = calculateDiscountForTenPercent(discount, product, offer, quantity, unitPrice)
	}
	if offer.offerType == FiveForAmount && quantityAsInt >= 5 {
		discount = calculateDiscountForFive(unitPrice, quantity, offer, numberOfXs, quantityAsInt, discount, product, x)
	}
	if discount != nil {
		receipt.addDiscount(*discount)
	}
}

func calculateDiscountForFive(unitPrice float64, quantity float64, offer SpecialOffer, numberOfXs int, quantityAsInt int, discount *Discount, product Product, x int) *Discount {
	var discountTotal = unitPrice*quantity - (offer.costAfterDiscount*float64(numberOfXs) + float64(quantityAsInt%5)*unitPrice)
	discount = &Discount{product: product, description: fmt.Sprintf("%d for %.2f", x, offer.costAfterDiscount), discountAmount: -discountTotal}
	return discount
}

func calculateDiscountForTenPercent(discount *Discount, product Product, offer SpecialOffer, quantity float64, unitPrice float64) *Discount {
	discount = &Discount{product: product, description: fmt.Sprintf("%.0f %% off", offer.costAfterDiscount), discountAmount: -quantity * unitPrice * offer.costAfterDiscount / 100.0}
	return discount
}

func calculateDiscountThreeForTwo(quantity float64, unitPrice float64, numberOfXs int, quantityAsInt int, discount *Discount, product Product) *Discount {
	var discountAmount = quantity*unitPrice - (float64(numberOfXs*2)*unitPrice + float64(quantityAsInt%3)*unitPrice)
	discount = &Discount{product: product, description: "3 for 2", discountAmount: -discountAmount}
	return discount
}

func calculateDiscountForTwo(offer SpecialOffer, quantityAsInt int, x int, unitPrice float64, quantity float64, discount *Discount, product Product) *Discount {
	var total = offer.costAfterDiscount*float64(quantityAsInt/x) + float64(quantityAsInt%2)*unitPrice
	var discountN = unitPrice*quantity - total
	discount = &Discount{product: product, description: fmt.Sprintf("2 for %.2f", offer.costAfterDiscount), discountAmount: -discountN}
	return discount
}
