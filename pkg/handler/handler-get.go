package handler

import (
	"net/http"
	"sort"
	"strconv"
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
}

var FirstRequestPath string

func (h *Handler) authMiddleware(c *gin.Context) {
	if v, err := c.Cookie("token"); err != nil {
		FirstRequestPath = c.Request.URL.Path
		c.Redirect(http.StatusFound, "/auth")

	} else {
		if v == h.nexusmanager.Config.AppPassword {
			c.Next()
		}
	}

}

func (h *Handler) getAuthRoot(c *gin.Context) {
	tmpl, _ := template.ParseFiles("template/auth.html")
	tmpl.Execute(c.Writer, nil)
}

func (h *Handler) postAuthRoot(c *gin.Context) {
	if password, ok := c.GetPostForm("password"); ok && password == h.nexusmanager.Config.AppPassword {
		username, _ := c.GetPostForm("username")
		logrus.Println("User " + username + " login success from IP " + c.Request.RemoteAddr)

		c.SetCookie("token", h.nexusmanager.Config.AppPassword, 3000, "/", c.Request.Host, false, true)
		c.Redirect(http.StatusFound, FirstRequestPath)
	} else {
		logrus.Println("Wrong password, IP ", "127.0.0.1", " password: ", password)
		c.Redirect(http.StatusFound, "/auth")
	}
}

func (h *Handler) getRoot(c *gin.Context) {
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
	tmpl.Execute(c.Writer, &DataStruct{list, *repo, c.Param("id")})
}
