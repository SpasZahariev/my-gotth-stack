package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type CheckoutSessionHandLer struct{}

func NewCheckoutSessionHandler() *CheckoutSessionHandLer {
	return &CheckoutSessionHandLer{}
}

func (h *CheckoutSessionHandLer) ServeHTTP(c *gin.Context) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Hot Sauce"),
					},
					UnitAmount: stripe.Int64(500), // Amount in cents
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String("https://yourdomain.com/success"),
		CancelURL:  stripe.String("https://yourdomain.com/cancel"),
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": s.ID})
}
