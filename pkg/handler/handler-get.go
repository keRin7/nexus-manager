package handler

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/keRin7/nexus-manager/nexusmanager"
	"github.com/sirupsen/logrus"
)

type TagsList struct {
	Data []string `json:"data"`
}

type Image struct {
	Id   string
	Size string
	Data string
}

type Images struct {
	Images []Image
}

type DataStruct struct {
	Images
	Repos       nexusmanager.Repositories
	CurrentRepo string
	Username    string
}

func (h *Handler) getAuthRoot(c *gin.Context) {
	tmpl, _ := template.ParseFiles("template/auth.html")
	tmpl.Execute(c.Writer, nil)
}

func (h *Handler) postAuthRoot(c *gin.Context) {
	if username, ok := c.GetPostForm("username"); ok {
		if password, ok := c.GetPostForm("password"); ok {
			if h.nexusmanager.Ldapcli.TryToBind(username, password) {
				logrus.Println("User " + username + " login success from IP " + c.Request.RemoteAddr)
				if token, err := h.nexusmanager.Auth.CreateToken(username); err != nil {
					logrus.Printf("Error generate JWT: %s", err.Error())
					c.Redirect(http.StatusFound, "/auth")
					return
				} else {
					c.SetCookie(authorizationHeader, "Bearer "+token, 3000, "/", "127.0.0.1:8080", false, true)
					c.Redirect(http.StatusFound, FirstRequestPath)
					return
				}
			}
		}
	}
	logrus.Println("Wrong password, IP ", c.Request.RemoteAddr)
	c.Redirect(http.StatusFound, "/auth")
	//if password, ok := c.GetPostForm("password"); ok && password == h.nexusmanager.Config.AppPassword {
	//	username, _ := c.GetPostForm("username")
	//	logrus.Println("User " + username + " login success from IP " + c.Request.RemoteAddr)

	//	c.SetCookie("token", h.nexusmanager.Config.AppPassword, 3000, "/", c.Request.Host, false, true)
	//	c.Redirect(http.StatusFound, FirstRequestPath)
	//} else {
	//	logrus.Println("Wrong password, IP ", "127.0.0.1", " password: ", password)
	//	c.Redirect(http.StatusFound, "/auth")
	//}
}

func (h *Handler) getRoot(c *gin.Context) {
	token, _ := h.GetToken(c)
	username := h.nexusmanager.Auth.GetUsername(token)
	logrus.Println("User: ", username, "get access to: ", h.nexusmanager.Config.Nexus_repo+"/"+c.Param("id"))

	var list Images
	tags := h.nexusmanager.ListTagsByImage(h.nexusmanager.Config.Nexus_repo + "/" + c.Param("id"))
	repo := h.nexusmanager.List()
	for _, v := range tags {
		data := h.nexusmanager.GetDataV1(h.nexusmanager.Config.Nexus_repo+"/"+c.Param("id"), v)
		size := h.nexusmanager.GetSize(h.nexusmanager.Config.Nexus_repo+"/"+c.Param("id"), v)
		list.Images = append(list.Images, Image{v, strconv.FormatInt(size/1024/1024, 10), data})
	}
	sort.Slice(list.Images, func(i, j int) bool {
		return list.Images[i].Data < list.Images[j].Data
	})
	tmpl, _ := template.ParseFiles("template/index.html")
	tmpl.Execute(c.Writer, &DataStruct{list, *repo, c.Param("id"), username})
}

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
