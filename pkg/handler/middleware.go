package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

var FirstRequestPath string

func (h *Handler) authMiddleware(c *gin.Context) {
	if token, err := h.GetToken(c); err != nil {
		FirstRequestPath = c.Request.URL.Path
		c.Redirect(http.StatusFound, "/auth")
	} else {
		if _, err := h.nexusmanager.Auth.ParseToken(token); err != nil {
			c.Redirect(http.StatusFound, "/auth")
		} else {
			c.Next()
		}
	}

}
