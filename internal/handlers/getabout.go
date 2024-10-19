package handlers

import (
	"goth/internal/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AboutHandLer struct{}

func NewAboutHandler() *AboutHandLer {
	return &AboutHandLer{}
}

func (h *AboutHandLer) ServeHTTP(c *gin.Context) {
	template := templates.About()
	err := templates.Layout(c, template, "My website").Render(c.Request.Context(), c.Writer)

	if err != nil {
		http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
