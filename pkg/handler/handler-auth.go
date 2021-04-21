package handler

import (
	"errors"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetToken(c *gin.Context) (string, error) {
	if header, err := c.Cookie(authorizationHeader); err != nil {
		return "", err
	} else {
		headers := strings.Split(header, " ")
		if len(headers) != 2 {
			return "", errors.New("auth header damaged")
		}
		return headers[1], nil
	}
}

func (h *Handler) getAuthRoot(c *gin.Context) {
	tmpl, _ := template.ParseFiles("template/auth.html")
	tmpl.Execute(c.Writer, nil)
}

func (h *Handler) postAuthRoot(c *gin.Context) {
	var username, password string
	var ok bool

	if username, ok = c.GetPostForm("username"); !ok {
		logrus.Println("POST username needed, IP ", c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, "/auth")
		return
	}
	if password, ok = c.GetPostForm("password"); !ok {
		logrus.Println("POST password needed, IP ", c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	if !h.nexusmanager.Ldapclient.TryToBind(username, password) {
		logrus.Println("Wrong username or password, IP ", c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, "/auth")
		return

	}

	logrus.Println("User " + username + " login success from IP " + c.Request.RemoteAddr)

	if token, err := h.nexusmanager.Auth.CreateToken(username); err != nil {
		logrus.Printf("Error generate JWT: %s", err.Error())
		c.Redirect(http.StatusFound, "/auth")
		return
	} else {
		c.SetCookie(authorizationHeader, "Bearer "+token, 3000, "/", "127.0.0.1:8080", false, true)
		c.Redirect(http.StatusFound, FirstRequestPath)
	}

}
