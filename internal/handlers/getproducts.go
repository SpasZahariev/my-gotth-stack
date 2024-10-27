package handlers

import (
	"goth/internal/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductsHandLer struct{}

func NewProductsHandler() *ProductsHandLer {
	return &ProductsHandLer{}
}

func (h *ProductsHandLer) ServeHTTP(c *gin.Context) {
	template := templates.Products()
	err := templates.Layout(c, template, "Products Page").Render(c.Request.Context(), c.Writer)

	if err != nil {
		http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
