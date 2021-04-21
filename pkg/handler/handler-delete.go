package handler

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) PostDelete(c *gin.Context) {
	var token string
	var err error
	var access bool
	var tag, image, repo string
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
				tag = elems[len(elems)-1]
				image = elems[len(elems)-2]
				repo = strings.Join(elems[:len(elems)-2], "/")
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
		tmpl, _ := template.ParseFiles("template/deleted.html")
		tmpl.Execute(c.Writer, gin.H{
			"Text": success_deleted_images,
			"Back": repo + "/" + image,
		})
		return

	}
	logrus.Println("Username:" + username + " tries to delete images ,but ... forbidden")

	tmpl, _ := template.ParseFiles("template/forbidden.html")
	tmpl.Execute(c.Writer, gin.H{"Text": "Forbidden"})

}
