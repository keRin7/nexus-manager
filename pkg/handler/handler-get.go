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

	if !h.nexusmanager.Ldapcli.TryToBind(username, password) {
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
	tmpl.Execute(c.Writer, &DataStruct{list, *repo, h.nexusmanager.Config.Nexus_repo + "/" + c.Param("id"), username})
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

func (h *Handler) PostDelete(c *gin.Context) {
	var token string
	var err error
	var access bool
	success_deleted_images := "Deleted images: <BR>"
	//Get JWT token from cookie, it need for getting username
	if token, err = h.GetToken(c); err != nil {
		c.Redirect(http.StatusFound, "/auth/")
	}
	//Getting username for more detail in out logs
	username := h.nexusmanager.Auth.GetUsername(token)

	if access, err = h.nexusmanager.Auth.CheckAccess(token); err != nil {
		c.Redirect(http.StatusFound, "/auth/")
	}

	if access {
		//Get and range all checkbox select ID
		for _, value := range c.PostFormArray("flexCheckChecked") {
			elems := strings.Split(value, "/")
			if len(elems) >= 3 {
				tag := elems[len(elems)-1]
				image := elems[len(elems)-2]
				repo := strings.Join(elems[:len(elems)-2], "/")
				if err := h.nexusmanager.DeleteImageByTag(repo+"/"+image, tag); err != nil {
					logrus.Println(err.Error())
				} else {
					success_deleted_images = success_deleted_images + "<BR>" + repo + "/" + image + ":" + tag
					logrus.Println("Endpoint: /delete,  Username:" + username + "image deleted success: +" + repo + "/" + image + tag)
				}
			} else {
				logrus.Println("Endpoint: /delete , Username:"+username+", The delete-url isn`t full:", strings.Join(elems, " "))
			}

		}
		tmpl, _ := template.ParseFiles("template/simple-text.html")
		tmpl.Execute(c.Writer, gin.H{"Text": success_deleted_images})
		return

	}
	logrus.Println("Username:" + username + " tries to delete images, forbidden")

	tmpl, _ := template.ParseFiles("template/simple-text.html")
	tmpl.Execute(c.Writer, gin.H{"Text": "Forbidden"})

}
