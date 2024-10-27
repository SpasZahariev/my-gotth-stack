package handlers

import (
	"encoding/json"
	"goth/internal/templates"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(c *gin.Context) {

	// user, ok := c.Request.Context().Value(middleware.UserKey).(*store.User)

	cookie, err := c.Cookie("userInfo")

	if err != nil {
		template := templates.GuestIndex()

		err := templates.Layout(c, template, "Home Page").Render(c.Request.Context(), c.Writer)

		if err != nil {
			http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
			return
		}

		return
	}

	// Decode the cookie value
	decodedValue, err := url.QueryUnescape(cookie)

	if err != nil {
		http.Error(c.Writer, "Error decoding cookie", http.StatusInternalServerError)
	}

	// Parse the JSON to get teh Email
	var userInfo struct {
		Email string `json:"primaryEmail"`
	}
	err = json.Unmarshal([]byte(decodedValue), &userInfo)
	if err != nil {
		http.Error(c.Writer, "Error parsing user info", http.StatusInternalServerError)
	}

	template := templates.Index(userInfo.Email)
	err = templates.Layout(c, template, "Logged In Home Page").Render(c.Request.Context(), c.Writer)

	if err != nil {
		http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
