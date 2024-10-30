package handlers

import (
	"goth/internal/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessfulPaymentHandLer struct{}

func NewSuccessfulPaymentHandler() *SuccessfulPaymentHandLer {
	return &SuccessfulPaymentHandLer{}
}

func (h *SuccessfulPaymentHandLer) ServeHTTP(c *gin.Context) {
	template := templates.SuccessfulPayment()
	err := templates.Layout(c, template, "Successful Payment Page!").Render(c.Request.Context(), c.Writer)

	if err != nil {
		http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
